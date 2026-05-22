import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import DetailChip from './DetailChip.vue';

describe('DetailChip', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(DetailChip, { props: { value: 'test' } });
    expect(wrapper.exists()).toBe(true);
  });

  it('displayLabel returns label when provided', () => {
    const wrapper = mountWithPlugins(DetailChip, {
      props: { value: 'val', label: 'My Label' },
    });
    expect(wrapper.vm.displayLabel).toBe('My Label');
  });

  it('displayLabel returns value as string when no label', () => {
    const wrapper = mountWithPlugins(DetailChip, {
      props: { value: 42 },
    });
    expect(wrapper.vm.displayLabel).toBe('42');
  });

  it('chipClasses includes size class', () => {
    const wrapper = mountWithPlugins(DetailChip, {
      props: { value: 'x', size: 'sm' },
    });
    expect(wrapper.vm.chipClasses).toContain('detail-chip--sm');
  });

  it('chipClasses includes dense when prop set', () => {
    const wrapper = mountWithPlugins(DetailChip, {
      props: { value: 'x', dense: true },
    });
    expect(wrapper.vm.chipClasses).toContain('detail-chip--dense');
  });
});
