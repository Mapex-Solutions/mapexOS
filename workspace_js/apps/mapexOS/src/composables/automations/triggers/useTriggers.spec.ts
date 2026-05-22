import { describe, it, expect, beforeEach } from 'vitest';
import { useTriggers } from './useTriggers';
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

/**
 * Helper to create a mock ConfiguredTrigger
 */
function makeTrigger(overrides: Partial<ConfiguredTrigger> = {}): ConfiguredTrigger {
  return {
    id: 'trig-1',
    templateId: 'tmpl-1',
    templateName: 'Temperature Alert',
    templateCategory: 'sensors',
    type: 'standard',
    variables: {},
    ...overrides,
  };
}

/**
 * Helper to create a mock TriggerGroup
 */
function makeTriggerGroup(overrides: Partial<TriggerGroup> = {}): TriggerGroup {
  return {
    id: 'grp-1',
    name: 'Group A',
    triggers: [makeTrigger()],
    ...overrides,
  };
}

describe('useTriggers', () => {
  beforeEach(() => {
    // Reset singleton state between tests
    const { individualTriggers, triggerGroups, externalTriggers } = useTriggers();
    individualTriggers.value = [];
    triggerGroups.value = [];
    externalTriggers.value = [];
  });

  describe('triggerOptions computed', () => {
    it('always includes "No Trigger" as first option', () => {
      const { triggerOptions } = useTriggers();
      expect(triggerOptions.value[0]).toEqual({
        label: 'No Trigger',
        value: 'No Trigger',
        icon: 'notifications_off',
      });
    });

    it('returns only "No Trigger" when no triggers exist', () => {
      const { triggerOptions } = useTriggers();
      expect(triggerOptions.value).toHaveLength(1);
    });

    it('includes single triggers with header when individual triggers exist', () => {
      const { triggerOptions, addIndividualTrigger } = useTriggers();

      addIndividualTrigger(makeTrigger({ id: 't-1', templateName: 'Temp Alert' }));

      const options = triggerOptions.value;
      // "No Trigger" + header + 1 trigger
      expect(options).toHaveLength(3);
      expect(options[1]!.value).toBe('single_triggers_header');
      expect(options[1]!.group).toBe('header');
      expect(options[2]!.value).toBe('t-1');
      expect(options[2]!.label).toBe('Temp Alert');
      expect(options[2]!.group).toBe('single_triggers');
    });

    it('includes trigger groups with header when groups exist', () => {
      const { triggerOptions, addTriggerGroup } = useTriggers();

      addTriggerGroup(makeTriggerGroup({ id: 'g-1', name: 'My Group' }));

      const options = triggerOptions.value;
      expect(options).toHaveLength(3);
      expect(options[1]!.value).toBe('trigger_groups_header');
      expect(options[2]!.value).toBe('g-1');
      expect(options[2]!.label).toBe('My Group');
    });

    it('skips groups with no triggers', () => {
      const { triggerOptions, addTriggerGroup } = useTriggers();

      addTriggerGroup(makeTriggerGroup({ id: 'g-1', name: 'Empty', triggers: [] }));

      // Only "No Trigger" + header (group is excluded because triggers.length === 0)
      expect(triggerOptions.value).toHaveLength(2);
    });

    it('includes external triggers with header when external triggers exist', () => {
      const { triggerOptions, addExternalTrigger } = useTriggers();

      addExternalTrigger(makeTrigger({ id: 'ext-1', templateName: 'Webhook', type: 'external' }));

      const options = triggerOptions.value;
      expect(options).toHaveLength(3);
      expect(options[1]!.value).toBe('external_triggers_header');
      expect(options[2]!.value).toBe('ext-1');
      expect(options[2]!.group).toBe('external_triggers');
    });

    it('includes all sections when all trigger types are present', () => {
      const { triggerOptions, addIndividualTrigger, addTriggerGroup, addExternalTrigger } = useTriggers();

      addIndividualTrigger(makeTrigger({ id: 't-1' }));
      addTriggerGroup(makeTriggerGroup({ id: 'g-1' }));
      addExternalTrigger(makeTrigger({ id: 'ext-1', type: 'external' }));

      // "No Trigger" + 3 headers + 3 items = 7
      expect(triggerOptions.value).toHaveLength(7);
    });
  });

  describe('individual trigger CRUD', () => {
    it('addIndividualTrigger adds to list', () => {
      const { individualTriggers, addIndividualTrigger } = useTriggers();
      const trigger = makeTrigger({ id: 't-1' });

      addIndividualTrigger(trigger);
      expect(individualTriggers.value).toHaveLength(1);
      expect(individualTriggers.value[0]!.id).toBe('t-1');
    });

    it('removeIndividualTrigger removes by id', () => {
      const { individualTriggers, addIndividualTrigger, removeIndividualTrigger } = useTriggers();

      addIndividualTrigger(makeTrigger({ id: 't-1' }));
      addIndividualTrigger(makeTrigger({ id: 't-2' }));

      removeIndividualTrigger('t-1');
      expect(individualTriggers.value).toHaveLength(1);
      expect(individualTriggers.value[0]!.id).toBe('t-2');
    });

    it('removeIndividualTrigger does nothing for unknown id', () => {
      const { individualTriggers, addIndividualTrigger, removeIndividualTrigger } = useTriggers();

      addIndividualTrigger(makeTrigger({ id: 't-1' }));
      removeIndividualTrigger('nonexistent');
      expect(individualTriggers.value).toHaveLength(1);
    });

    it('updateIndividualTrigger replaces trigger at index', () => {
      const { individualTriggers, addIndividualTrigger, updateIndividualTrigger } = useTriggers();

      addIndividualTrigger(makeTrigger({ id: 't-1', templateName: 'Old' }));

      const updated = makeTrigger({ id: 't-1', templateName: 'New' });
      updateIndividualTrigger('t-1', updated);

      expect(individualTriggers.value[0]!.templateName).toBe('New');
    });

    it('updateIndividualTrigger does nothing for unknown id', () => {
      const { individualTriggers, addIndividualTrigger, updateIndividualTrigger } = useTriggers();

      addIndividualTrigger(makeTrigger({ id: 't-1', templateName: 'Original' }));
      updateIndividualTrigger('nonexistent', makeTrigger({ id: 'x', templateName: 'X' }));

      expect(individualTriggers.value[0]!.templateName).toBe('Original');
    });
  });

  describe('external trigger CRUD', () => {
    it('addExternalTrigger adds to list', () => {
      const { externalTriggers, addExternalTrigger } = useTriggers();

      addExternalTrigger(makeTrigger({ id: 'ext-1', type: 'external' }));
      expect(externalTriggers.value).toHaveLength(1);
    });

    it('removeExternalTrigger removes by id', () => {
      const { externalTriggers, addExternalTrigger, removeExternalTrigger } = useTriggers();

      addExternalTrigger(makeTrigger({ id: 'ext-1' }));
      addExternalTrigger(makeTrigger({ id: 'ext-2' }));

      removeExternalTrigger('ext-1');
      expect(externalTriggers.value).toHaveLength(1);
      expect(externalTriggers.value[0]!.id).toBe('ext-2');
    });

    it('updateExternalTrigger replaces trigger at index', () => {
      const { externalTriggers, addExternalTrigger, updateExternalTrigger } = useTriggers();

      addExternalTrigger(makeTrigger({ id: 'ext-1', templateName: 'Old' }));
      updateExternalTrigger('ext-1', makeTrigger({ id: 'ext-1', templateName: 'Updated' }));

      expect(externalTriggers.value[0]!.templateName).toBe('Updated');
    });
  });

  describe('trigger group CRUD', () => {
    it('addTriggerGroup adds to list', () => {
      const { triggerGroups, addTriggerGroup } = useTriggers();

      addTriggerGroup(makeTriggerGroup({ id: 'g-1' }));
      expect(triggerGroups.value).toHaveLength(1);
    });

    it('removeTriggerGroup removes by id', () => {
      const { triggerGroups, addTriggerGroup, removeTriggerGroup } = useTriggers();

      addTriggerGroup(makeTriggerGroup({ id: 'g-1' }));
      addTriggerGroup(makeTriggerGroup({ id: 'g-2' }));

      removeTriggerGroup('g-1');
      expect(triggerGroups.value).toHaveLength(1);
      expect(triggerGroups.value[0]!.id).toBe('g-2');
    });

    it('updateTriggerGroup replaces group at index', () => {
      const { triggerGroups, addTriggerGroup, updateTriggerGroup } = useTriggers();

      addTriggerGroup(makeTriggerGroup({ id: 'g-1', name: 'Old' }));
      updateTriggerGroup('g-1', makeTriggerGroup({ id: 'g-1', name: 'New' }));

      expect(triggerGroups.value[0]!.name).toBe('New');
    });

    it('updateTriggerGroup does nothing for unknown id', () => {
      const { triggerGroups, addTriggerGroup, updateTriggerGroup } = useTriggers();

      addTriggerGroup(makeTriggerGroup({ id: 'g-1', name: 'Original' }));
      updateTriggerGroup('nonexistent', makeTriggerGroup({ id: 'x', name: 'X' }));

      expect(triggerGroups.value[0]!.name).toBe('Original');
    });
  });

  describe('singleton behavior', () => {
    it('shares state across multiple useTriggers calls', () => {
      const instance1 = useTriggers();
      const instance2 = useTriggers();

      instance1.addIndividualTrigger(makeTrigger({ id: 'shared-1' }));

      expect(instance2.individualTriggers.value).toHaveLength(1);
      expect(instance2.individualTriggers.value[0]!.id).toBe('shared-1');
    });
  });
});
