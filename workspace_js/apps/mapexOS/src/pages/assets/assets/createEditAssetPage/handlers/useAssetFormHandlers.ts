import type { Ref } from 'vue';
import type { QForm } from 'quasar';
import type { AssetFormData, AssetFormState, StepMetaEntry } from '../interfaces';
import type { StepperVerticalItem } from '@components/steppers';
import type { AssetTemplateResponse, RouteGroupResponse } from '@mapexos/schemas';

import { computed } from 'vue';
import { buildStepperTree } from '@components/steppers';

import { apis } from '@services/mapex';
import { notifySuccess } from '@utils/alert/notify';
import { handleApiError } from '@utils/error';
import { useAddAssetTranslations } from '@src/composables/i18n/pages/assets/addAsset/useAddAssetTranslations';
import { useOrganizationStore } from '@stores/organization';
import { useLogger } from '@composables/useLogger';

import { STEP } from '../constants';

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

  // Visual groups for the vertical stepper. Grouping is UI-only — the leaves are
  // the real navigable steps; the payload is unaffected. A group with no visible
  // leaf (e.g. dataRouting for a gateway) simply does not render.
  const setupGroup = { id: 'setup', label: t.steps.groups.setup, icon: 'mdi-cog-outline' };
  const identificationGroup = { id: 'identification', label: t.steps.groups.identification, icon: 'mdi-fingerprint' };
  const connectivityGroup = { id: 'connectivity', label: t.steps.groups.connectivity, icon: 'mdi-wifi' };
  const routingGroup = { id: 'routing', label: t.steps.groups.routing, icon: 'mdi-routes' };
  const finalizationGroup = { id: 'finalization', label: t.steps.groups.finalization, icon: 'mdi-clipboard-check' };

  // Ordered step metadata (leaves). The visible list is derived from this by kind;
  // a LoRaWAN gateway drops the assetTemplate + routeGroups steps. `group` places a
  // leaf under a stepper group header without changing navigation or payload.
  const stepMeta: StepMetaEntry[] = [
    { id: STEP.TYPE, icon: 'mdi-shape-outline', label: t.steps.stepType.label, description: t.steps.stepType.description, group: setupGroup },
    { id: STEP.ASSET_TEMPLATE, icon: 'mdi-file-document', label: t.steps.step2.label, description: t.steps.step2.description, group: setupGroup },
    { id: STEP.IDENTIFICATION, icon: 'mdi-fingerprint', label: t.steps.step1.leafLabel, description: t.steps.step1.leafDescription, group: identificationGroup },
    { id: STEP.ATTRIBUTES, icon: 'mdi-tag-multiple', label: t.steps.step1.attributesLabel, description: t.steps.step1.attributesDescription, group: identificationGroup },
    { id: STEP.CONNECTIVITY, icon: 'mdi-wifi', label: t.steps.step4.label, description: t.steps.step4.description, group: connectivityGroup },
    { id: STEP.HEALTH_MONITORING, icon: 'mdi-heart-pulse', label: t.steps.step5.label, description: t.steps.step5.description, group: connectivityGroup },
    { id: STEP.ROUTE_GROUPS, icon: 'mdi-routes', label: t.steps.step3.label, description: t.steps.step3.description, group: routingGroup },
    { id: STEP.REVIEW, icon: 'mdi-clipboard-check', label: t.steps.step6.label, description: t.steps.step6.description, group: finalizationGroup },
  ];

  // A LoRaWAN gateway is radio infrastructure with no data model and no routing,
  // so it skips the asset-template and route-group steps.
  const isGatewaySelected = computed(
    () => assetData.value.protocol === 'LORAWAN' && assetData.value.lorawanConfig?.kind === 'gateway'
  );

  /** Flat visible leaves for the current selection — drives navigation + rendering. */
  const visibleSteps = computed(() =>
    stepMeta
      .filter((s) => !(isGatewaySelected.value && (s.id === STEP.ASSET_TEMPLATE || s.id === STEP.ROUTE_GROUPS)))
      .map((s) => ({ id: s.id, icon: s.icon, title: s.label.value, description: s.description.value, group: s.group }))
  );

  /** Nested tree for the vertical stepper (navigation still runs on visibleSteps).
   * Group labels are resolved from their computed refs before grouping. */
  const stepperTree = computed<StepperVerticalItem[]>(() =>
    buildStepperTree(
      visibleSteps.value.map((leaf) => ({
        title: leaf.title,
        description: leaf.description,
        icon: leaf.icon,
        ...(leaf.group ? { group: { id: leaf.group.id, label: leaf.group.label.value, icon: leaf.group.icon } } : {}),
      }))
    )
  );

  /** Id of the step at the current 1-based position. */
  const currentStepId = computed(() => visibleSteps.value[currentStep.value - 1]?.id);

  // Per-step form refs, keyed by step id (the Type/Health/Review steps have none).
  const formRefByStep: Partial<Record<string, Ref<QForm | null>>> = {
    [STEP.IDENTIFICATION]: step1FormRef,
    [STEP.ASSET_TEMPLATE]: step2FormRef,
    [STEP.ROUTE_GROUPS]: step3FormRef,
    [STEP.CONNECTIVITY]: step4FormRef,
  };

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
   * Check if the Next button should be disabled for the current step. The
   * template and route-group steps gate on their selection; other steps let
   * their own form validation handle it.
   */
  const isNextButtonDisabled = computed(() => {
    if (currentStepId.value === STEP.ASSET_TEMPLATE) {
      return !assetData.value.assetTemplateId;
    }
    if (currentStepId.value === STEP.ROUTE_GROUPS) {
      return assetData.value.routeGroupIds.length === 0;
    }
    return false;
  });

  /**
   * Validate the current step's form (when one exists) before moving forward,
   * then move to the clamped target position in the visible step list.
   * @param {number} step - Target 1-based position in the visible step list.
   */
  async function changeStep(step: number) {
    const target = Math.min(Math.max(step, 1), visibleSteps.value.length);
    if (target > currentStep.value) {
      const formRef = formRefByStep[currentStepId.value ?? ''];
      if (formRef?.value) {
        const valid = await formRef.value.validate();
        if (!valid) return;
      }
    }
    currentStep.value = target;
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
  /**
   * Builds the protocol.lorawan block from the form. The `kind` selects a
   * disjoint shape: a gateway (frequency plan + connection auth) or a device
   * (identity + activation keys). Identity/profile fields always travel; secret
   * key material (OTAA appKey/nwkKey or the ABP session keys) is sent only when
   * the operator typed it — blank on edit keeps the existing envelope-encrypted
   * keys, mirroring the MQTT password path. The backend seals the keys and never
   * returns them.
   */
  function buildLorawanConfig(): Record<string, unknown> {
    const lw = assetData.value.lorawanConfig;

    if (lw.kind === 'gateway') {
      const gw: Record<string, unknown> = {
        authMode: lw.gateway.authMode,
        frequencyPlanId: lw.gateway.frequencyPlanId,
      };
      if (lw.gateway.frequencyPlanIds.length) gw.frequencyPlanIds = lw.gateway.frequencyPlanIds;
      // Latitude/longitude are the asset-level location (entered once on the
      // form); only altitude is gateway-specific. Fold them into the gateway
      // block the projection/LNS reads.
      if (assetData.value.latitude != null) gw.latitude = assetData.value.latitude;
      if (assetData.value.longitude != null) gw.longitude = assetData.value.longitude;
      if (lw.gateway.authMode === 'cert' && lw.gateway.certTTL) {
        gw.certTTL = { value: lw.gateway.certTTL.value, unit: lw.gateway.certTTL.unit };
      }
      // Key mode: send the token only when the operator filled it (blank on
      // edit keeps the existing hash, mirroring the device-key handling).
      if (lw.gateway.authMode === 'key' && lw.gateway.apiKey) {
        gw.apiKey = lw.gateway.apiKey;
      }
      return { kind: 'gateway', gateway: gw };
    }

    const device: Record<string, unknown> = {
      kind: 'device',
      // The asset UUID (Step 1) IS the device EUI — collected once, not a
      // separate field, so they can never diverge.
      devEui: assetData.value.assetId,
      region: lw.region,
      class: lw.class,
      macVersion: lw.macVersion,
      phyVersion: lw.phyVersion,
      activation: lw.activation,
    };
    if (lw.activation === 'otaa') {
      device.joinEui = lw.joinEui;
      if (lw.appKey) device.appKey = lw.appKey;
      if (lw.nwkKey) device.nwkKey = lw.nwkKey;
    } else {
      if (lw.devAddr) device.devAddr = lw.devAddr;
      if (lw.nwkSKey) device.nwkSKey = lw.nwkSKey;
      if (lw.appSKey) device.appSKey = lw.appSKey;
    }
    return device;
  }

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
    } else if (protocolType === 'lorawan') {
      protocolConfig.lorawan = buildLorawanConfig();
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
      attributes: assetData.value.attributes.map((a) => ({
        label: a.label,
        kind: a.kind,
        value: a.value,
        searchable: a.searchable ?? true,
      })),
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
    visibleSteps,
    stepperTree,
    currentStepId,
    isNextButtonDisabled,
    changeStep,
    handleStepChange,
    submitForm,
  };
}
