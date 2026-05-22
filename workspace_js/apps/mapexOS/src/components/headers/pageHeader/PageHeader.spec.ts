import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import PageHeader from './PageHeader.vue';
import type { PageHeaderProps } from './interfaces';

vi.mock('@utils/translation', () => ({
  useTS: () => (key: string) => key,
}));

vi.mock('@components/tooltips', () => ({
  AppTooltip: { template: '<span />' },
}));

vi.mock('@components/dialogs/infoModal', () => ({
  InfoModal: { template: '<div />', props: ['modelValue', 'icon', 'title', 'description', 'items', 'docsUrl', 'docsLabel'] },
}));

const baseProps: PageHeaderProps = {
  title: 'Test Page',
};

describe('PageHeader', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('rendering', () => {
    it('should render with minimal props (title only)', () => {
      const wrapper = mountWithPlugins(PageHeader, { props: baseProps });
      expect(wrapper.exists()).toBe(true);
      expect(wrapper.text()).toContain('Test Page');
    });

    it('should render the icon when provided', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: { ...baseProps, icon: 'dashboard', iconColor: 'accent' },
      });
      const icon = wrapper.find('q-icon-stub');
      expect(icon.exists()).toBe(true);
      expect(icon.attributes('name')).toBe('dashboard');
      expect(icon.attributes('color')).toBe('accent');
    });

    it('should use default icon color "primary" when iconColor is not provided', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: { ...baseProps, icon: 'dashboard' },
      });
      const icon = wrapper.find('q-icon-stub');
      expect(icon.attributes('color')).toBe('primary');
    });

    it('should render description when provided', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: { ...baseProps, description: 'A description' },
      });
      expect(wrapper.text()).toContain('A description');
    });

    it('should not render description element when not provided', () => {
      const wrapper = mountWithPlugins(PageHeader, { props: baseProps });
      const descDiv = wrapper.find('.text-subtitle1.text-grey-7');
      expect(descDiv.exists()).toBe(false);
    });
  });

  describe('button', () => {
    it('should render the action button when button prop is provided', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: {
          ...baseProps,
          button: { label: 'Create', icon: 'add', color: 'positive', id: 'create-btn' },
        },
      });
      const btn = wrapper.find('q-btn-stub[id="create-btn"]');
      expect(btn.exists()).toBe(true);
      expect(btn.attributes('label')).toBe('Create');
      expect(btn.attributes('icon')).toBe('add');
    });

    it('should not render the action button when button prop is not provided', () => {
      const wrapper = mountWithPlugins(PageHeader, { props: baseProps });
      const btns = wrapper.findAll('q-btn-stub');
      // No buttons at all when no button, no tour, no info
      expect(btns.length).toBe(0);
    });

    it('should apply default button styles', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: {
          ...baseProps,
          button: { label: 'Go' },
        },
      });
      const btn = wrapper.find('q-btn-stub[label="Go"]');
      expect(btn.attributes('color')).toBe('grey-7');
    });
  });

  describe('tour button', () => {
    it('should render tour button when tour.enabled is true', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: { ...baseProps, tour: { enabled: true } },
      });
      const tourBtn = wrapper.find('#tour-start-btn');
      expect(tourBtn.exists()).toBe(true);
    });

    it('should emit "start-tour" when tour button is clicked', async () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: { ...baseProps, tour: { enabled: true } },
      });
      await wrapper.find('#tour-start-btn').trigger('click');
      expect(wrapper.emitted('start-tour')).toHaveLength(1);
    });

    it('should not render tour button when tour is not provided', () => {
      const wrapper = mountWithPlugins(PageHeader, { props: baseProps });
      expect(wrapper.find('#tour-start-btn').exists()).toBe(false);
    });
  });

  describe('info button and modal', () => {
    const infoConfig = {
      title: 'Info Title',
      description: 'Info description',
      items: [{ text: 'Feature 1' }],
      docsUrl: 'https://docs.example.com',
    };

    it('should render info button when info prop is provided and tour is not', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: { ...baseProps, info: infoConfig },
      });
      const infoBtn = wrapper.find('.info-button');
      expect(infoBtn.exists()).toBe(true);
    });

    it('should NOT render info button when tour is enabled (tour takes priority)', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: { ...baseProps, info: infoConfig, tour: { enabled: true } },
      });
      expect(wrapper.find('.info-button').exists()).toBe(false);
      expect(wrapper.find('#tour-start-btn').exists()).toBe(true);
    });

    it('should open info modal when info button is clicked', async () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: { ...baseProps, info: infoConfig },
      });
      expect(wrapper.vm.infoModalOpen).toBe(false);
      await wrapper.find('.info-button').trigger('click');
      expect(wrapper.vm.infoModalOpen).toBe(true);
    });
  });

  describe('slots', () => {
    it('should render header-extra slot content', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: baseProps,
        slots: { 'header-extra': '<span class="extra-content">Extra</span>' },
      });
      expect(wrapper.find('.extra-content').exists()).toBe(true);
    });

    it('should render actions slot content', () => {
      const wrapper = mountWithPlugins(PageHeader, {
        props: baseProps,
        slots: { actions: '<button class="custom-action">Action</button>' },
      });
      expect(wrapper.find('.custom-action').exists()).toBe(true);
    });
  });

  describe('computed', () => {
    it('should compute infoTooltipText from i18n key', () => {
      const wrapper = mountWithPlugins(PageHeader, { props: baseProps });
      expect(wrapper.vm.infoTooltipText).toBe(
        'components.headers.pageHeader.infoTooltip',
      );
    });
  });
});
