import { ref, computed } from 'vue';

interface TriggerOption {
  label: string;
  value: string;
  icon?: string;
  group?: string;
}

interface ConfiguredTrigger {
  id: string;
  templateName: string;
  [key: string]: any;
}

interface TriggerGroup {
  id: string;
  name: string;
  triggers: ConfiguredTrigger[];
}

// Create singleton state
const individualTriggers = ref<ConfiguredTrigger[]>([]);
const triggerGroups = ref<TriggerGroup[]>([]);
const externalTriggers = ref<ConfiguredTrigger[]>([]);

export function useTriggers() {
  const triggerOptions = computed<TriggerOption[]>(() => {
    const options: TriggerOption[] = [
      {
        label: 'No Trigger',
        value: 'No Trigger',
        icon: 'notifications_off'
      }
    ];

    // Add individual triggers under "Single Triggers" group with header
    if (individualTriggers.value.length > 0) {
      options.push({
        label: '⚡ Single Triggers',
        value: 'single_triggers_header',
        group: 'header'
      });

      individualTriggers.value.forEach(trigger => {
        options.push({
          label: `${trigger.templateName}`,
          value: trigger.id,
          group: 'single_triggers',
          icon: 'bolt'
        });
      });
    }

    // Add group triggers under "Trigger Groups" with header
    if (triggerGroups.value.length > 0) {
      options.push({
        label: '👥 Trigger Groups',
        value: 'trigger_groups_header',
        group: 'header'
      });

      triggerGroups.value.forEach(group => {
        if (group.triggers.length > 0) {
          options.push({
            label: `${group.name}`,
            value: group.id,
            group: 'trigger_groups',
            icon: 'group'
          });
        }
      });
    }

    // Add external triggers under "External Triggers" with header
    if (externalTriggers.value.length > 0) {
      options.push({
        label: '🔗 External Triggers',
        value: 'external_triggers_header',
        group: 'header'
      });

      externalTriggers.value.forEach(trigger => {
        options.push({
          label: `${trigger.templateName}`,
          value: trigger.id,
          group: 'external_triggers',
          icon: 'link'
        });
      });
    }

    return options;
  });

  function addIndividualTrigger(trigger: ConfiguredTrigger): void {
    individualTriggers.value.push(trigger);
  }

  function removeIndividualTrigger(triggerId: string): void {
    const index = individualTriggers.value.findIndex(t => t.id === triggerId);
    if (index !== -1) {
      individualTriggers.value.splice(index, 1);
    }
  }

  function updateIndividualTrigger(triggerId: string, trigger: ConfiguredTrigger): void {
    const index = individualTriggers.value.findIndex(t => t.id === triggerId);
    if (index !== -1) {
      individualTriggers.value[index] = trigger;
    }
  }

  function addExternalTrigger(trigger: ConfiguredTrigger): void {
    externalTriggers.value.push(trigger);
  }

  function removeExternalTrigger(triggerId: string): void {
    const index = externalTriggers.value.findIndex(t => t.id === triggerId);
    if (index !== -1) {
      externalTriggers.value.splice(index, 1);
    }
  }

  function updateExternalTrigger(triggerId: string, trigger: ConfiguredTrigger): void {
    const index = externalTriggers.value.findIndex(t => t.id === triggerId);
    if (index !== -1) {
      externalTriggers.value[index] = trigger;
    }
  }

  function addTriggerGroup(group: TriggerGroup): void {
    triggerGroups.value.push(group);
  }

  function removeTriggerGroup(groupId: string): void {
    const index = triggerGroups.value.findIndex(g => g.id === groupId);
    if (index !== -1) {
      triggerGroups.value.splice(index, 1);
    }
  }

  function updateTriggerGroup(groupId: string, group: TriggerGroup): void {
    const index = triggerGroups.value.findIndex(g => g.id === groupId);
    if (index !== -1) {
      triggerGroups.value[index] = group;
    }
  }

  return {
    individualTriggers,
    triggerGroups,
    externalTriggers,
    triggerOptions,
    addIndividualTrigger,
    removeIndividualTrigger,
    updateIndividualTrigger,
    addExternalTrigger,
    removeExternalTrigger,
    updateExternalTrigger,
    addTriggerGroup,
    removeTriggerGroup,
    updateTriggerGroup
  };
}
