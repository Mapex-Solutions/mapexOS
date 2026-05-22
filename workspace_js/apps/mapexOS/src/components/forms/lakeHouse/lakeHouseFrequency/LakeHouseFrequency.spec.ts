import { describe, it, expect, vi } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import LakeHouseFrequency from './LakeHouseFrequency.vue';

vi.mock('./constants', () => ({
  DEFAULT_FREQUENCY_CONFIG: { type: 'day', interval: 1, time: '09:00' },
}));

const BASE_MODEL = {
  frequency: { type: 'day' as const, interval: 1, time: '09:00' },
};

describe('LakeHouseFrequency', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('computes showInterval as true for day', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.showInterval).toBe(true);
  });

  it('computes showInterval as false for week', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: { frequency: { type: 'week', weekdays: [], time: '09:00' } } },
    });
    expect(wrapper.vm.showInterval).toBe(false);
  });

  it('computes showTime as true for day', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.showTime).toBe(true);
  });

  it('computes showTime as false for minute', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: { frequency: { type: 'minute', interval: 5 } } },
    });
    expect(wrapper.vm.showTime).toBe(false);
  });

  it('computes showWeekdays as true for week', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: { frequency: { type: 'week', weekdays: [], time: '09:00' } } },
    });
    expect(wrapper.vm.showWeekdays).toBe(true);
  });

  it('computes showDayOfMonth as true for month', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: { frequency: { type: 'month', dayOfMonth: 1, time: '09:00' } } },
    });
    expect(wrapper.vm.showDayOfMonth).toBe(true);
  });

  it('computes intervalSuffix correctly for day', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.intervalSuffix).toBe('day');
  });

  it('computes intervalSuffix plural for multiple days', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: { frequency: { type: 'day', interval: 3, time: '09:00' } } },
    });
    expect(wrapper.vm.intervalSuffix).toBe('days');
  });

  it('computes frequencyDescription for daily', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.frequencyDescription).toContain('Daily');
  });

  it('validates time format correctly', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.isValidTime('14:30')).toBe(true);
    expect(wrapper.vm.isValidTime('25:00')).toBe(false);
    expect(wrapper.vm.isValidTime('')).toBe(true);
  });

  it('selectFrequency updates model to minute', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.selectFrequency('minute');
    expect(wrapper.vm.modelRef.frequency.type).toBe('minute');
    expect(wrapper.vm.modelRef.frequency.interval).toBe(1);
  });

  it('selectFrequency updates model to week with empty weekdays', () => {
    const wrapper = mountWithPlugins(LakeHouseFrequency, {
      props: { modelValue: BASE_MODEL },
    });
    wrapper.vm.selectFrequency('week');
    expect(wrapper.vm.modelRef.frequency.type).toBe('week');
    expect(wrapper.vm.modelRef.frequency.weekdays).toEqual([]);
  });
});
