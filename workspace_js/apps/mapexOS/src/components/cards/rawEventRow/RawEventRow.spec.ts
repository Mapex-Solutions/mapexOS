import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import RawEventRow from './RawEventRow.vue';
import type { RawEventProps } from './interfaces';

const baseEvent: RawEventProps = {
  id: 'raw-1',
  asset: {
    name: 'Sensor A',
    description: 'Temperature sensor',
    icon: 'thermostat',
    type: 'sensor',
  },
  type: 'telemetry',
  status: 'high',
  protocol: 'MQTT',
  created: '2025-01-15T14:00:00Z',
  values: { temperature: 42 },
};

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(RawEventRow, {
    props: { event: baseEvent, ...overrides },
  });
}

describe('RawEventRow', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('getCardClass returns class based on status', () => {
    const wrapper = factory();
    expect(wrapper.vm.getCardClass()).toBe('event-card--high');
  });

  it('getBorderClass returns border class based on status', () => {
    const wrapper = factory();
    expect(wrapper.vm.getBorderClass()).toBe('event-card__border--high');
  });

  it('getStatusColor returns red-6 for high status', () => {
    const wrapper = factory();
    expect(wrapper.vm.getStatusColor()).toBe('red-6');
  });

  it('getStatusColor returns orange-6 for medium status', () => {
    const wrapper = factory({
      event: { ...baseEvent, status: 'medium' },
    });
    expect(wrapper.vm.getStatusColor()).toBe('orange-6');
  });

  it('getStatusColor returns green-6 for low status', () => {
    const wrapper = factory({
      event: { ...baseEvent, status: 'low' },
    });
    expect(wrapper.vm.getStatusColor()).toBe('green-6');
  });

  it('getAssetIconColor returns red-5 for high status', () => {
    const wrapper = factory();
    expect(wrapper.vm.getAssetIconColor()).toBe('red-5');
  });

  it('getProtocolColor returns purple-6 for MQTT', () => {
    const wrapper = factory();
    expect(wrapper.vm.getProtocolColor()).toBe('purple-6');
  });

  it('getProtocolColor returns blue-6 for HTTP', () => {
    const wrapper = factory({
      event: { ...baseEvent, protocol: 'HTTP' },
    });
    expect(wrapper.vm.getProtocolColor()).toBe('blue-6');
  });

  it('getProtocolColor returns grey-6 for unknown protocol', () => {
    const wrapper = factory({
      event: { ...baseEvent, protocol: 'CUSTOM' },
    });
    expect(wrapper.vm.getProtocolColor()).toBe('grey-6');
  });

  it('formatDate returns a formatted string', () => {
    const wrapper = factory();
    const formatted = wrapper.vm.formatDate('2025-01-15T14:00:00Z');
    expect(typeof formatted).toBe('string');
    expect(formatted.length).toBeGreaterThan(0);
  });
});
