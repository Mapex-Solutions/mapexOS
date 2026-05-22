import type { Ref } from 'vue';
import type { QForm } from 'quasar';
import type { AssetFormData, AssetFormState } from '../interfaces';
import type { AssetTemplateResponse, RouteGroupResponse } from '@mapexos/schemas';

import { computed } from 'vue';

import { apis } from '@services/mapex';
import { notifySuccess } from '@utils/alert/notify';
import { handleApiError } from '@utils/error';
import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';
import { useOrganizationStore } from '@stores/organization';
import { useLogger } from '@composables/useLogger';

/**
 * Result returned by submitForm. Page-level orchestration (post-save
 * cert dialog vs. straight navigation) reads this so the wizard can
 * branch between cert-mode MQTT creates (open download dialog) and
 * the rest (navigate immediately). The handler no longer navigates
 * on its own.
 */
export interface SubmitFormResult {
  ok: boolean;
  isCreate: boolean;
  isMqtt: boolean;
  isCert: boolean;
  assetUUID?: string;
  assetName?: string;
}

const logger = useLogger('useAssetFormHandlers');

interface UseAssetFormHandlersParams {
  assetData: Ref<AssetFormData>;
  formState: Ref<AssetFormState>;
  currentStep: Ref<number>;
  isEditMode: Ref<boolean>;
  assetId: Ref<string | undefined>;
  isSaving: Ref<boolean>;
  step1FormRef: Ref<QForm | null>;
  step2FormRef: Ref<QForm | null>;
  step3FormRef: Ref<QForm | null>;
  step4FormRef: Ref<QForm | null>;
}

export function useAssetFormHandlers(params: UseAssetFormHandlersParams) {
  const {
    assetData,
    formState,
    currentStep,
    isEditMode,
    assetId,
    isSaving,
    step1FormRef,
    step2FormRef,
    step3FormRef,
    step4FormRef,
  } = params;

  const t = useAddAssetTranslations();
  const orgStore = useOrganizationStore();

  /**
   * Handle template selection
   */
  function onTemplateSelected(template: AssetTemplateResponse | null) {
    logger.debug('onTemplateSelected called with:', template);
    formState.value.selectedTemplate = template;
  }

  /**
   * Handle route groups selection
   */
  function onRouteGroupsSelected(routeGroups: RouteGroupResponse[]) {
    formState.value.selectedRouteGroups = routeGroups;
  }

  /**
   * Check if Next button should be disabled
   */
  const isNextButtonDisabled = computed(() => {
    if (currentStep.value === 1) {
      return false; // Let validate() handle it
    }
    if (currentStep.value === 2) {
      return !assetData.value.assetTemplateId;
    }
    if (currentStep.value === 3) {
      return assetData.value.routeGroupIds.length === 0;
    }
    if (currentStep.value === 4) {
      return false; // Let validate() handle it
    }
    return false;
  });

  /**
   * Validate and change step
   */
  async function changeStep(step: number) {
    if (currentStep.value === 1 && step > 1 && step1FormRef.value) {
      const valid = await step1FormRef.value.validate();
      if (!valid) return;
    }
    if (currentStep.value === 2 && step > 2 && step2FormRef.value) {
      const valid = await step2FormRef.value.validate();
      if (!valid) return;
    }
    if (currentStep.value === 3 && step > 3 && step3FormRef.value) {
      const valid = await step3FormRef.value.validate();
      if (!valid) return;
    }
    if (currentStep.value === 4 && step > 4 && step4FormRef.value) {
      const valid = await step4FormRef.value.validate();
      if (!valid) return;
    }

    currentStep.value = step;
  }

  /**
   * Wrapper for step navigation (non-async)
   */
  function handleStepChange(step: number): void {
    void changeStep(step);
  }

  /**
   * Builds the request payload from the current form state.
   *
   * MQTT password handling:
   *   - CREATE: password is REQUIRED and flows verbatim as
   *     `protocol.mqtt.password` (plaintext). The backend bcrypts it
   *     before persisting; the plaintext never returns on responses.
   *   - EDIT: password is OPTIONAL. We OMIT the field when blank so a
   *     PATCH that only touches other fields keeps the existing
   *     bcrypt hash. When provided, the backend rotates the hash
   *     using the same cost as the create path.
   */
  function buildAssetPayload(): any {
    const protocolType = assetData.value.protocol.toLowerCase() as 'http' | 'mqtt' | 'lorawan';

    const protocolConfig: any = { type: protocolType };
    if (protocolType === 'http') {
      protocolConfig.http = {};
    } else if (protocolType === 'mqtt') {
      const mqtt: {
        clientId: string;
        username: string;
        authType: string;
        password?: string;
        certTTL?: { value: number; unit: string };
      } = {
        clientId: assetData.value.mqttConfig.clientId,
        username: assetData.value.mqttConfig.username,
        authType: assetData.value.mqttConfig.authType,
      };
      // Only send password when the asset is in password mode AND the
      // operator typed a value (blank on edit = keep existing hash).
      const pwd = assetData.value.mqttConfig.password;
      if (assetData.value.mqttConfig.authType === 'password' && pwd && pwd.length > 0) {
        mqtt.password = pwd;
      }
      // certTTL travels only in cert mode — password-mode assets never
      // issue a device cert so persisting the field there would be
      // misleading + the broker contract rejects it anyway.
      if (assetData.value.mqttConfig.authType === 'cert' && assetData.value.mqttConfig.certTTL) {
        mqtt.certTTL = {
          value: assetData.value.mqttConfig.certTTL.value,
          unit: assetData.value.mqttConfig.certTTL.unit,
        };
      }
      protocolConfig.mqtt = mqtt;
    }

    const healthMonitor = assetData.value.healthMonitor.enabled
      ? {
          enabled: true,
          thresholdMinutes: assetData.value.healthMonitor.thresholdMinutes,
          requiredMisses: assetData.value.healthMonitor.requiredMisses,
          heartbeatMode: assetData.value.healthMonitor.heartbeatMode ?? 'implicit',
          offlineRouteGroupIds: assetData.value.healthMonitor.offlineRouteGroupIds,
          onlineRouteGroupIds: assetData.value.healthMonitor.onlineRouteGroupIds,
        }
      : { enabled: false };

    const payload: any = {
      name: assetData.value.name,
      enabled: assetData.value.enabled,
      debugEnabled: assetData.value.debugEnabled,
      description: assetData.value.description,
      assetUUID: assetData.value.assetId,
      assetTemplateId: assetData.value.assetTemplateId!,
      routeGroupIds: assetData.value.routeGroupIds,
      protocol: protocolConfig,
      latitude: assetData.value.latitude ?? undefined,
      longitude: assetData.value.longitude ?? undefined,
      healthMonitor,
    };

    if (!isEditMode.value) {
      payload.orgId = orgStore.selectedOrganizationId!;
    }

    return payload;
  }

  /**
   * Submit form — POST on create, PATCH on edit. Returns the outcome
   * so the page can decide what happens next (e.g. open the cert
   * download dialog for newly-created MQTT assets, navigate back to
   * the list for everything else). Drops `mqttConfig.password` from
   * local state on success so the field re-renders blank instead of
   * revealing the plaintext on the next view. On error, the handler
   * surfaces a notification and returns ok=false; the page is then
   * free to stay on the form.
   */
  async function submitForm(): Promise<SubmitFormResult> {
    isSaving.value = true;

    const isCreate = !(isEditMode.value && assetId.value);
    const isMqtt = assetData.value.protocol === 'MQTT';
    const isCert = isMqtt && assetData.value.mqttConfig.authType === 'cert';

    try {
      const payload = buildAssetPayload();

      let createdUUID: string | undefined;
      if (!isCreate) {
        logger.debug('Updating Asset');
        await apis.assets.asset.update({ assetId: assetId.value! }, payload);
        notifySuccess({ message: t.notifications.updated.value, timeout: 3000 });
      } else {
        logger.debug('Creating Asset');
        const created = await apis.assets.asset.create(payload);
        createdUUID = created?.assetUUID ?? assetData.value.assetId;
        notifySuccess({ message: t.notifications.created.value, timeout: 3000 });
      }

      assetData.value.mqttConfig.password = '';

      const finalUUID = isCreate ? createdUUID : assetId.value;
      const result: SubmitFormResult = { ok: true, isCreate, isMqtt, isCert };
      if (finalUUID) result.assetUUID = finalUUID;
      if (assetData.value.name) result.assetName = assetData.value.name;
      return result;
    } catch (error: any) {
      handleApiError(error, {
        customMessages: {
          409: t.notifications.alreadyExists,
          422: t.notifications.validationFailed,
          network: t.notifications.networkError,
        },
        defaultMessage: isEditMode.value
          ? t.notifications.updateFailed.value
          : t.notifications.creationFailed,
        timeout: 5000,
      });
      return { ok: false, isCreate, isMqtt, isCert };
    } finally {
      isSaving.value = false;
    }
  }

  return {
    onTemplateSelected,
    onRouteGroupsSelected,
    isNextButtonDisabled,
    changeStep,
    handleStepChange,
    submitForm,
  };
}
