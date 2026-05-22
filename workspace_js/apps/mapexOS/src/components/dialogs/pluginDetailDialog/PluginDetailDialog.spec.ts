import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import PluginDetailDialog from './PluginDetailDialog.vue';

describe('PluginDetailDialog', () => {
  const defaultProps = {
    modelValue: true,
    name: 'Telegram',
    author: 'MapexOS',
    version: '1.0.0',
    description: 'Send messages via Telegram',
    brandIconUrl: '',
    icon: 'send',
    color: '#0088cc',
    categoryLabel: 'Messaging',
    tags: ['chat', 'notification'],
    loading: false,
    nodeTypes: [
      {
        type: 'telegram/sendMessage',
        label: 'Send Message',
        icon: 'send',
        color: '#0088cc',
        description: 'Send a message to a Telegram chat',
        inputCount: 1,
        outputCount: 1,
      },
    ],
    installed: false,
    installing: false,
    installDisabled: false,
    installLabel: 'Install',
    installingLabel: 'Installing...',
    installedLabel: 'Installed',
    nodeTypesLabel: 'Node Types',
    loadingLabel: 'Loading...',
    inputsLabel: 'inputs',
    outputsLabel: 'outputs',
  };

  it('renders without errors', () => {
    const wrapper = mountWithPlugins(PluginDetailDialog, { props: defaultProps });
    expect(wrapper.exists()).toBe(true);
  });

  it('syncs isOpen ref with modelValue prop', () => {
    const wrapper = mountWithPlugins(PluginDetailDialog, { props: defaultProps });
    expect(wrapper.vm.isOpen).toBe(true);
  });

  it('sets isOpen to false when modelValue is false', () => {
    const wrapper = mountWithPlugins(PluginDetailDialog, {
      props: { ...defaultProps, modelValue: false },
    });
    expect(wrapper.vm.isOpen).toBe(false);
  });

  it('emits update:modelValue when isOpen changes', async () => {
    const wrapper = mountWithPlugins(PluginDetailDialog, { props: defaultProps });
    wrapper.vm.isOpen = false;
    await wrapper.vm.$nextTick();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('emits update:modelValue(false) on handleClose', () => {
    const wrapper = mountWithPlugins(PluginDetailDialog, { props: defaultProps });
    (wrapper.vm).handleClose();
    expect(wrapper.emitted('update:modelValue')).toBeTruthy();
    expect(wrapper.emitted('update:modelValue')![0]).toEqual([false]);
  });

  it('accepts all required props', () => {
    const wrapper = mountWithPlugins(PluginDetailDialog, { props: defaultProps });
    expect(wrapper.props('name')).toBe('Telegram');
    expect(wrapper.props('author')).toBe('MapexOS');
    expect(wrapper.props('version')).toBe('1.0.0');
    expect(wrapper.props('description')).toBe('Send messages via Telegram');
    expect(wrapper.props('categoryLabel')).toBe('Messaging');
  });

  it('receives tags as array', () => {
    const wrapper = mountWithPlugins(PluginDetailDialog, { props: defaultProps });
    expect(wrapper.props('tags')).toEqual(['chat', 'notification']);
  });

  it('receives nodeTypes array', () => {
    const wrapper = mountWithPlugins(PluginDetailDialog, { props: defaultProps });
    expect(wrapper.props('nodeTypes')).toHaveLength(1);
    expect(wrapper.props('nodeTypes')[0].type).toBe('telegram/sendMessage');
  });

  it('emits install event', () => {
    const wrapper = mountWithPlugins(PluginDetailDialog, { props: defaultProps });
    wrapper.vm.$emit('install');
    expect(wrapper.emitted('install')).toBeTruthy();
  });
});
