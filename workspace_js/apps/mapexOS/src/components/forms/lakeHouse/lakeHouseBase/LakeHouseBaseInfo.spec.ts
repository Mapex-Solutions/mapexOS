import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import LakeHouseBaseInfo from './LakeHouseBaseInfo.vue';

const BASE_MODEL = {
  name: 'Test Lake',
  status: true,
  description: 'A description',
};

describe('LakeHouseBaseInfo', () => {
  it('renders without errors', () => {
    const wrapper = mountWithPlugins(LakeHouseBaseInfo, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('receives modelValue with correct name', () => {
    const wrapper = mountWithPlugins(LakeHouseBaseInfo, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.modelRef.name).toBe('Test Lake');
  });

  it('receives modelValue with correct status', () => {
    const wrapper = mountWithPlugins(LakeHouseBaseInfo, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.modelRef.status).toBe(true);
  });

  it('receives modelValue with correct description', () => {
    const wrapper = mountWithPlugins(LakeHouseBaseInfo, {
      props: { modelValue: BASE_MODEL },
    });
    expect(wrapper.vm.modelRef.description).toBe('A description');
  });
});
