import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import AssetRow from './AssetRow.vue';
import type { AssetRowProps } from './interfaces';

const baseAsset: AssetRowProps = {
  id: '1',
  name: 'Temperature Sensor',
  enabled: true,
  assetUUID: 'uuid-abc-123',
  assetTypeName: 'Sensor',
  protocol: { type: 'mqtt' },
};

function factory(assetOverrides: Partial<AssetRowProps> = {}) {
  return mountWithPlugins(AssetRow, {
    props: {
      asset: { ...baseAsset, ...assetOverrides },
    },
  });
}

describe('AssetRow', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('expanded starts as false', () => {
    const wrapper = factory();
    expect(wrapper.vm.expanded).toBe(false);
  });

  it('toggleExpand flips expanded state and emits click', () => {
    const wrapper = factory();
    wrapper.vm.toggleExpand();
    expect(wrapper.vm.expanded).toBe(true);
    expect(wrapper.emitted('click')).toBeTruthy();
  });

  it('toggleExpand does not emit click when collapsing', () => {
    const wrapper = factory();
    wrapper.vm.toggleExpand(); // expand
    wrapper.vm.toggleExpand(); // collapse
    expect(wrapper.vm.expanded).toBe(false);
    // click only emitted once (on expand)
    expect(wrapper.emitted('click')).toHaveLength(1);
  });

  it('getIconColor returns primary when enabled', () => {
    const wrapper = factory({ enabled: true });
    expect(wrapper.vm.getIconColor()).toBe('primary');
  });

  it('getIconColor returns grey-5 when disabled', () => {
    const wrapper = factory({ enabled: false });
    expect(wrapper.vm.getIconColor()).toBe('grey-5');
  });

  it('getTypeColorForChip returns blue when enabled', () => {
    const wrapper = factory({ enabled: true });
    expect(wrapper.vm.getTypeColorForChip()).toBe('blue');
  });

  it('getTypeColorForChip returns grey when disabled', () => {
    const wrapper = factory({ enabled: false });
    expect(wrapper.vm.getTypeColorForChip()).toBe('grey');
  });

  it('getProtocolColorForChip returns purple for mqtt', () => {
    const wrapper = factory({ protocol: { type: 'mqtt' } });
    expect(wrapper.vm.getProtocolColorForChip()).toBe('purple');
  });

  it('getProtocolColorForChip returns blue for http', () => {
    const wrapper = factory({ protocol: { type: 'http' } });
    expect(wrapper.vm.getProtocolColorForChip()).toBe('blue');
  });

  it('getProtocolColorForChip returns orange for lorawan', () => {
    const wrapper = factory({ protocol: { type: 'lorawan' } });
    expect(wrapper.vm.getProtocolColorForChip()).toBe('orange');
  });

  it('getProtocolColorForChip returns grey when protocol is absent', () => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { protocol: _omit, ...baseAssetWithoutProtocol } = baseAsset;
    const wrapper = mountWithPlugins(AssetRow, {
      props: { asset: baseAssetWithoutProtocol },
    });
    expect(wrapper.vm.getProtocolColorForChip()).toBe('grey');
  });
});
