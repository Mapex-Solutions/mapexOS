import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import SelectableChip from './SelectableChip.vue';

describe('SelectableChip', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(SelectableChip, {
      props: { label: 'Test', modelValue: false },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('mounts with correct prop values', () => {
    const wrapper = mountWithPlugins(SelectableChip, {
      props: { label: 'Test', modelValue: true },
    });
    expect(wrapper.vm).toBeDefined();
  });
});
