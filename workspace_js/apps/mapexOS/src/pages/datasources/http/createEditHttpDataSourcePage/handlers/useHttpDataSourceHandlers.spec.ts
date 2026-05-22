import { describe, it, expect, beforeEach, vi } from 'vitest';
import { ref } from 'vue';
import type { Ref } from 'vue';

import type { HttpDataSource } from '../interfaces';

import { useHttpDataSourceHandlers } from './useHttpDataSourceHandlers';

/** Mock notify utilities */
vi.mock('@utils/alert/notify', () => ({
  notifySuccess: vi.fn(),
  notifyFail: vi.fn(),
}));

/** Mock logger */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    debug: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  }),
}));

import { notifySuccess, notifyFail } from '@utils/alert/notify';

/**
 * Factory for creating a default HttpDataSource ref
 */
function createDataSource(overrides: Partial<HttpDataSource> = {}): HttpDataSource {
  return {
    name: 'Test DS',
    description: '',
    enabled: true,
    mode: null,
    protocol: null,
    enableWorkingHours: false,
    daysOfWeek: [],
    timeIntervals: [{ startTime: '09:00', endTime: '17:00' }],
    timezone: 'UTC',
    enableRateLimit: false,
    rateLimitType: null,
    rateLimitValue: 0,
    burstCapacity: 0,
    actionOnExceed: null,
    authType: null,
    apiKey: { headerApiKey: '', valueApiKey: '' },
    jwt: { secretKey: '', headerName: '' },
    ipWhitelist: { addresses: [] },
    oauth2: { jwksUrl: '' },
    bindingMode: null,
    directAssetId: null,
    directAssetIdPath: null,
    assetTemplateIds: [],
    customUuidPaths: [{ path: '' }],
    payloadExample: '',
    ...overrides,
  };
}

describe('useHttpDataSourceHandlers', () => {
  let dataSource: ReturnType<typeof ref<HttpDataSource>>;

  beforeEach(() => {
    vi.clearAllMocks();
    dataSource = ref(createDataSource());
  });

  function setup() {
    return useHttpDataSourceHandlers(dataSource as Ref<HttpDataSource>);
  }

  describe('addInterval', () => {
    it('adds a default time interval', () => {
      const { addInterval } = setup();

      addInterval();

      expect(dataSource.value!.timeIntervals).toHaveLength(2);
      expect(dataSource.value!.timeIntervals[1]).toEqual({ startTime: '09:00', endTime: '17:00' });
    });
  });

  describe('removeInterval', () => {
    it('removes interval at the specified index', () => {
      dataSource.value!.timeIntervals.push({ startTime: '10:00', endTime: '18:00' });
      const { removeInterval } = setup();

      removeInterval(0);

      expect(dataSource.value!.timeIntervals).toHaveLength(1);
      expect(dataSource.value!.timeIntervals[0]).toEqual({ startTime: '10:00', endTime: '18:00' });
    });

    it('does not remove the last remaining interval', () => {
      const { removeInterval } = setup();

      removeInterval(0);

      expect(dataSource.value!.timeIntervals).toHaveLength(1);
    });
  });

  describe('addMapping', () => {
    it('adds an empty custom UUID path', () => {
      const { addMapping } = setup();

      addMapping();

      expect(dataSource.value!.customUuidPaths).toHaveLength(2);
      expect(dataSource.value!.customUuidPaths[1]).toEqual({ path: '' });
    });
  });

  describe('removeMapping', () => {
    it('removes mapping at the specified index', () => {
      dataSource.value!.customUuidPaths.push({ path: 'device.id' });
      const { removeMapping } = setup();

      removeMapping(0);

      expect(dataSource.value!.customUuidPaths).toHaveLength(1);
      expect(dataSource.value!.customUuidPaths[0]).toEqual({ path: 'device.id' });
    });

    it('does not remove the last remaining mapping', () => {
      const { removeMapping } = setup();

      removeMapping(0);

      expect(dataSource.value!.customUuidPaths).toHaveLength(1);
    });
  });

  describe('testMapping', () => {
    it('shows success notification when paths extract values', () => {
      dataSource.value!.payloadExample = JSON.stringify({ device: { uuid: '123' } });
      dataSource.value!.finalUuidPaths = ['device.uuid'];
      const { testMapping } = setup();

      testMapping();

      expect(notifySuccess).toHaveBeenCalledWith(
        expect.objectContaining({ message: expect.stringContaining('device.uuid: 123') }),
      );
    });

    it('shows fail notification when no paths configured', () => {
      dataSource.value!.payloadExample = '{}';
      dataSource.value!.finalUuidPaths = [];
      const { testMapping } = setup();

      testMapping();

      expect(notifyFail).toHaveBeenCalledWith(
        expect.objectContaining({ message: 'No UUID paths configured to test' }),
      );
    });

    it('shows fail notification on invalid JSON', () => {
      dataSource.value!.payloadExample = 'not-json';
      const { testMapping } = setup();

      testMapping();

      expect(notifyFail).toHaveBeenCalledWith(
        expect.objectContaining({ message: 'Invalid JSON payload or path format' }),
      );
    });
  });

  describe('saveDataSource', () => {
    it('returns true and notifies on success', () => {
      const { saveDataSource } = setup();

      const result = saveDataSource();

      expect(result).toBe(true);
      expect(notifySuccess).toHaveBeenCalledWith(
        expect.objectContaining({ message: 'Data Source created successfully!' }),
      );
    });
  });
});
