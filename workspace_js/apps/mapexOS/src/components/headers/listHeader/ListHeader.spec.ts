import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import ListHeader from './ListHeader.vue';
import type { ListHeadersProps } from './interfaces';

const baseProps: ListHeadersProps = {
  title: 'Asset List',
  icon: 'devices',
  button: {
    label: 'Create New',
    icon: 'add',
  },
};

describe('ListHeader', () => {
  describe('rendering', () => {
    it('should render with required props', () => {
      const wrapper = mountWithPlugins(ListHeader, { props: baseProps });
      expect(wrapper.exists()).toBe(true);
      expect(wrapper.text()).toContain('Asset List');
    });

    it('should render the header icon', () => {
      const wrapper = mountWithPlugins(ListHeader, { props: baseProps });
      const icon = wrapper.find('q-icon-stub');
      expect(icon.exists()).toBe(true);
      expect(icon.attributes('name')).toBe('devices');
      expect(icon.attributes('color')).toBe('primary');
    });

    it('should render the action button with correct label and icon', () => {
      const wrapper = mountWithPlugins(ListHeader, { props: baseProps });
      const btn = wrapper.find('q-btn-stub');
      expect(btn.exists()).toBe(true);
      expect(btn.attributes('label')).toBe('Create New');
      expect(btn.attributes('icon')).toBe('add');
    });
  });

  describe('button props', () => {
    it('should use default color "primary" when button.color is not provided', () => {
      const wrapper = mountWithPlugins(ListHeader, { props: baseProps });
      const btn = wrapper.find('q-btn-stub');
      expect(btn.attributes('color')).toBe('primary');
    });

    it('should use custom color when button.color is provided', () => {
      const wrapper = mountWithPlugins(ListHeader, {
        props: {
          ...baseProps,
          button: { ...baseProps.button, color: 'negative' },
        },
      });
      const btn = wrapper.find('q-btn-stub');
      expect(btn.attributes('color')).toBe('negative');
    });

    it('should pass the "to" route to the button', () => {
      const wrapper = mountWithPlugins(ListHeader, {
        props: {
          ...baseProps,
          button: { ...baseProps.button, to: '/assets/create' },
        },
      });
      const btn = wrapper.find('q-btn-stub');
      expect(btn.attributes('to')).toBe('/assets/create');
    });
  });

  describe('different titles and icons', () => {
    it('should render any title passed as prop', () => {
      const wrapper = mountWithPlugins(ListHeader, {
        props: { ...baseProps, title: 'Workflows' },
      });
      expect(wrapper.text()).toContain('Workflows');
    });

    it('should render any icon passed as prop', () => {
      const wrapper = mountWithPlugins(ListHeader, {
        props: { ...baseProps, icon: 'account_tree' },
      });
      const icon = wrapper.find('q-icon-stub');
      expect(icon.attributes('name')).toBe('account_tree');
    });
  });
});
