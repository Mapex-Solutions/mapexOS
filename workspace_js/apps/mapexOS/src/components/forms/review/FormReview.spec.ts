import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import FormReview from './FormReview.vue';

/** Mock components */
vi.mock('@components/chips', () => ({
  DetailChip: { template: '<span />' },
}));

describe('FormReview', () => {
  const sampleSections = [
    {
      stepNumber: 1,
      label: 'Basic Info',
      icon: { name: 'info', color: 'primary' },
      testId: 'section-basic',
      fields: [
        { label: 'Name', value: 'Test Asset', type: 'text' as const },
        { label: 'Status', value: 'Active', type: 'badge' as const, badgeColors: 'green-6' },
      ],
    },
    {
      stepNumber: 2,
      label: 'Configuration',
      icon: { name: 'settings' },
      testId: 'section-config',
      fields: [
        { label: 'Enabled', value: true, type: 'boolean' as const },
        { label: 'Created', value: '2024-01-15T10:30:00Z', type: 'datetime' as const, format: 'date' },
      ],
    },
  ];

  it('renders without errors with sections', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('renders without errors with empty sections', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: [] },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('handleEditSection emits editSection with step number', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    wrapper.vm.handleEditSection(2);
    expect(wrapper.emitted('editSection')).toBeTruthy();
    expect(wrapper.emitted('editSection')![0]![0]).toBe(2);
  });

  it('getTextValue returns dash for null value', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    const field = { label: 'Test', value: null, type: 'text' as const };
    expect(wrapper.vm.getTextValue(field)).toBe('—');
  });

  it('getTextValue returns dash for empty string', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    const field = { label: 'Test', value: '', type: 'text' as const };
    expect(wrapper.vm.getTextValue(field)).toBe('—');
  });

  it('getTextValue returns stringified value for string', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    const field = { label: 'Test', value: 'Hello', type: 'text' as const };
    expect(wrapper.vm.getTextValue(field)).toBe('Hello');
  });

  it('getTextValue returns stringified value for number', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    const field = { label: 'Test', value: 42, type: 'text' as const };
    expect(wrapper.vm.getTextValue(field)).toBe('42');
  });

  it('getBooleanLabel returns Active for truthy value', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    expect(wrapper.vm.getBooleanLabel(true)).toBe('Active');
  });

  it('getBooleanLabel returns Inactive for falsy value', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    expect(wrapper.vm.getBooleanLabel(false)).toBe('Inactive');
  });

  it('getBooleanColor returns green-6 for truthy', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    expect(wrapper.vm.getBooleanColor(true)).toBe('green-6');
  });

  it('getBooleanColor returns grey-6 for falsy', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    expect(wrapper.vm.getBooleanColor(false)).toBe('grey-6');
  });

  it('getColor returns string badgeColors directly', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    const field = { label: 'Status', value: 'ok', type: 'badge' as const, badgeColors: 'green-6' };
    expect(wrapper.vm.getColor(field)).toBe('green-6');
  });

  it('getColor returns mapped color from object badgeColors', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    const field = {
      label: 'Status',
      value: 'active',
      type: 'badge' as const,
      badgeColors: { active: 'green-6', inactive: 'grey-6' },
    };
    expect(wrapper.vm.getColor(field)).toBe('green-6');
  });

  it('getColor returns primary when no badgeColors', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    const field = { label: 'Status', value: 'ok', type: 'badge' as const };
    expect(wrapper.vm.getColor(field)).toBe('primary');
  });

  it('formatJson returns formatted JSON string', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    expect(wrapper.vm.formatJson({ a: 1 })).toBe(JSON.stringify({ a: 1 }, null, 2));
  });

  it('formatJson returns dash for null', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    expect(wrapper.vm.formatJson(null)).toBe('—');
  });

  it('formatDateTime returns dash for empty value', () => {
    const wrapper = mountWithPlugins(FormReview, {
      props: { sections: sampleSections },
    });
    expect(wrapper.vm.formatDateTime('')).toBe('—');
  });
});
