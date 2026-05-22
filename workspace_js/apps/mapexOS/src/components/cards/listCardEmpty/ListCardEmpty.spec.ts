import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ListCardEmpty from './ListCardEmpty.vue';

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(ListCardEmpty, {
    props: {
      icon: 'inbox',
      title: 'No items found',
      description: 'Try adjusting your filters',
      ...overrides,
    },
  });
}

describe('ListCardEmpty', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('receives icon prop correctly', () => {
    const wrapper = factory({ icon: 'search_off' });
    expect(wrapper.props('icon')).toBe('search_off');
  });

  it('receives title prop correctly', () => {
    const wrapper = factory({ title: 'Empty State' });
    expect(wrapper.props('title')).toBe('Empty State');
  });

  it('receives description prop correctly', () => {
    const wrapper = factory({ description: 'No data available' });
    expect(wrapper.props('description')).toBe('No data available');
  });

  it('receives optional buttonLabel prop', () => {
    const wrapper = factory({ buttonLabel: 'Add New' });
    expect(wrapper.props('buttonLabel')).toBe('Add New');
  });

  it('receives optional buttonIcon prop', () => {
    const wrapper = factory({ buttonIcon: 'add' });
    expect(wrapper.props('buttonIcon')).toBe('add');
  });

  it('emits button-click event', () => {
    const wrapper = factory();
    wrapper.vm.handleClick();
    expect(wrapper.emitted('button-click')).toBeTruthy();
    expect(wrapper.emitted('button-click')).toHaveLength(1);
  });
});
