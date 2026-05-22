import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ListPagination from './ListPagination.vue';

const mockQuasar = {
  screen: {
    lt: { sm: false, md: false, lg: false },
  },
};

vi.mock('quasar', () => ({
  useQuasar: () => mockQuasar,
}));

describe('ListPagination', () => {
  beforeEach(() => {
    mockQuasar.screen.lt.sm = false;
    mockQuasar.screen.lt.md = false;
    mockQuasar.screen.lt.lg = false;
  });

  it('renders pagination when totalPages > 1', () => {
    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 5 },
    });

    expect(wrapper.find('q-pagination-stub').exists()).toBe(true);
  });

  it('does not render when totalPages is 1', () => {
    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 1 },
    });

    expect(wrapper.find('q-pagination-stub').exists()).toBe(false);
  });

  it('does not render when totalPages is 0', () => {
    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 0 },
    });

    expect(wrapper.find('q-pagination-stub').exists()).toBe(false);
  });

  it('passes modelValue to q-pagination', () => {
    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 3, totalPages: 10 },
    });

    expect(wrapper.find('q-pagination-stub').attributes('model-value')).toBe('3');
  });

  it('passes totalPages as max to q-pagination', () => {
    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 10 },
    });

    expect(wrapper.find('q-pagination-stub').attributes('max')).toBe('10');
  });

  it('passes default color props', () => {
    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 5 },
    });

    const pagination = wrapper.find('q-pagination-stub');
    expect(pagination.attributes('color')).toBe('primary');
    expect(pagination.attributes('active-color')).toBe('primary');
  });

  it('passes custom color props', () => {
    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 5, color: 'accent', activeColor: 'secondary' },
    });

    const pagination = wrapper.find('q-pagination-stub');
    expect(pagination.attributes('color')).toBe('accent');
    expect(pagination.attributes('active-color')).toBe('secondary');
  });

  it('emits update:modelValue and change on page change', async () => {
    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 5 },
    });

    const pagination = wrapper.findComponent({ name: 'q-pagination' });
    await pagination.vm.$emit('update:model-value', 3);

    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([3]);
    expect(wrapper.emitted('change')).toBeTruthy();
    expect(wrapper.emitted('change')![0]).toEqual([3]);
  });

  it('computes maxPages as 10 for desktop (default, no lt flags)', () => {
    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 20 },
    });

    expect(wrapper.find('q-pagination-stub').attributes('max-pages')).toBe('10');
  });

  it('computes maxPages as 3 for mobile screens', () => {
    mockQuasar.screen.lt.sm = true;
    mockQuasar.screen.lt.md = true;
    mockQuasar.screen.lt.lg = true;

    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 20 },
    });

    expect(wrapper.find('q-pagination-stub').attributes('max-pages')).toBe('3');
  });

  it('computes maxPages as 5 for tablet screens', () => {
    mockQuasar.screen.lt.sm = false;
    mockQuasar.screen.lt.md = true;
    mockQuasar.screen.lt.lg = true;

    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 20 },
    });

    expect(wrapper.find('q-pagination-stub').attributes('max-pages')).toBe('5');
  });

  it('computes maxPages as 7 for laptop screens', () => {
    mockQuasar.screen.lt.sm = false;
    mockQuasar.screen.lt.md = false;
    mockQuasar.screen.lt.lg = true;

    const wrapper = mountWithPlugins(ListPagination, {
      props: { modelValue: 1, totalPages: 20 },
    });

    expect(wrapper.find('q-pagination-stub').attributes('max-pages')).toBe('7');
  });
});
