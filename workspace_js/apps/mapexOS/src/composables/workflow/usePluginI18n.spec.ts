import { describe, it, expect, vi, beforeEach } from 'vitest';

/** Mock vue-i18n's useI18n */
const mockGlobalT = vi.fn((key: string) => `translated:${key}`);

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: mockGlobalT }),
}));

import { usePluginI18n } from './usePluginI18n';

describe('usePluginI18n', () => {
  beforeEach(() => {
    mockGlobalT.mockClear();
    mockGlobalT.mockImplementation((key: string) => `translated:${key}`);
  });

  describe('key resolution', () => {
    it('prefixes keys with wf.{pluginId}.', () => {
      const { t } = usePluginI18n('core-flow-control');
      t('nodes.end.config.terminateWithError');

      expect(mockGlobalT).toHaveBeenCalledWith(
        'wf.core-flow-control.nodes.end.config.terminateWithError',
      );
    });

    it('returns the translated string', () => {
      const { t } = usePluginI18n('my-plugin');
      const result = t('label');

      expect(result).toBe('translated:wf.my-plugin.label');
    });
  });

  describe('interpolation params', () => {
    it('passes params to the global t function', () => {
      const { t } = usePluginI18n('core-http');
      const params = { count: 3, name: 'test' };

      t('messages.fetched', params);

      expect(mockGlobalT).toHaveBeenCalledWith(
        'wf.core-http.messages.fetched',
        params,
      );
    });

    it('does not pass params when undefined', () => {
      const { t } = usePluginI18n('core-http');
      t('simple.key');

      expect(mockGlobalT).toHaveBeenCalledWith('wf.core-http.simple.key');
      // Ensure it was called with exactly one argument
      expect(mockGlobalT.mock.calls[0]).toHaveLength(1);
    });
  });

  describe('different plugin IDs', () => {
    it('scopes to the correct plugin namespace', () => {
      const plugin1 = usePluginI18n('plugin-a');
      const plugin2 = usePluginI18n('plugin-b');

      plugin1.t('key');
      plugin2.t('key');

      expect(mockGlobalT).toHaveBeenCalledWith('wf.plugin-a.key');
      expect(mockGlobalT).toHaveBeenCalledWith('wf.plugin-b.key');
    });
  });
});
