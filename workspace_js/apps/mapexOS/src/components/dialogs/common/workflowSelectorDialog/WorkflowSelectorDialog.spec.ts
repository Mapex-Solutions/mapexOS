import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import WorkflowSelectorDialog from './WorkflowSelectorDialog.vue';

vi.mock('@utils/translation', () => ({
  useTS: () => (key: string) => key,
}));

vi.mock('@components/dialogs/common/genericSelectorDialog', () => ({
  GenericSelectorDialog: { name: 'GenericSelectorDialog', template: '<div />' },
}));

vi.mock('./constants', () => ({
  DUMMY_WORKFLOWS: [
    { id: 'wf-1', name: 'Workflow A', description: 'Desc A', enabled: true },
    { id: 'wf-2', name: 'Workflow B', description: 'Desc B', enabled: false },
    { id: 'wf-3', name: 'Workflow C', description: 'Desc C', enabled: true },
  ],
}));

const BASE_PROPS = {
  modelValue: false,
};

describe('WorkflowSelectorDialog', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.exists()).toBe(true);
  });

  it('starts with empty workflows', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.workflows).toEqual([]);
  });

  it('starts with loading as false', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.loading).toBe(false);
  });

  it('computes selectedIds from selectedWorkflowId', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, {
      props: { ...BASE_PROPS, selectedWorkflowId: 'wf-1' },
    });
    expect(wrapper.vm.selectedIds).toEqual(['wf-1']);
  });

  it('computes selectedIds as empty when no selectedWorkflowId', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.selectedIds).toEqual([]);
  });

  it('computes statusOptions with three entries', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.statusOptions).toHaveLength(3);
  });

  it('handles search query update', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, { props: BASE_PROPS });
    wrapper.vm.handleSearch('workflow');
    expect(wrapper.vm.searchQuery).toBe('workflow');
  });

  it('handles select by emitting workflow', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, { props: BASE_PROPS });
    const workflow = { id: 'wf-1', name: 'Workflow A', description: 'Desc A', enabled: true };
    wrapper.vm.handleSelect([workflow]);
    const emitted = wrapper.emitted('select')!;
    expect(emitted[0]![0]).toEqual(workflow);
  });

  it('does not emit on empty select', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, { props: BASE_PROPS });
    wrapper.vm.handleSelect([]);
    expect(wrapper.emitted('select')).toBeFalsy();
  });

  it('computes filteredWorkflows with empty list by default', () => {
    const wrapper = mountWithPlugins(WorkflowSelectorDialog, { props: BASE_PROPS });
    expect(wrapper.vm.filteredWorkflows).toEqual([]);
  });
});
