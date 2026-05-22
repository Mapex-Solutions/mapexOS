import { describe, it, expect, beforeEach, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useAuthStore } from './index';

/** Mock useLogger — not globally mocked in setup.ts */
vi.mock('@composables/useLogger', () => ({
  useLogger: () => ({
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

/** Mock storage utility */
vi.mock('@utils/storages', () => ({
  default: {
    local: {
      set: vi.fn(),
      get: vi.fn(),
      remove: vi.fn(),
    },
    session: {
      set: vi.fn(),
      get: vi.fn(),
      remove: vi.fn(),
    },
  },
}));

/** Shared mock instances so assertions work across module boundaries */
const mockOrgStore = {
  initializeAfterLogin: vi.fn().mockResolvedValue(undefined),
  clearCoverage: vi.fn(),
};

const mockPermStore = {
  clearPermissions: vi.fn(),
};

/** Mock organization store — auth actions depend on it */
vi.mock('@stores/organization', () => ({
  useOrganizationStore: () => mockOrgStore,
}));

/** Mock permission store — logout clears permissions */
vi.mock('@stores/permission', () => ({
  usePermissionStore: () => mockPermStore,
}));

describe('AuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    vi.clearAllMocks();
  });

  // ────────────────────────────────────────────────────────────────────────
  // State
  // ────────────────────────────────────────────────────────────────────────

  describe('state', () => {
    it('has default email', () => {
      const store = useAuthStore();
      expect(store.email).toBe('admin@mapex.global');
    });

    it('has default password', () => {
      const store = useAuthStore();
      expect(store.password).toBe('mapex123');
    });

    it('has keepConnected true by default', () => {
      const store = useAuthStore();
      expect(store.keepConnected).toBe(true);
    });

    it('has empty accessToken by default', () => {
      const store = useAuthStore();
      expect(store.accessToken).toBe('');
    });

    it('has empty refreshToken by default', () => {
      const store = useAuthStore();
      expect(store.refreshToken).toBe('');
    });

    it('has user null by default', () => {
      const store = useAuthStore();
      expect(store.user).toBeNull();
    });

    it('has loading false by default', () => {
      const store = useAuthStore();
      expect(store.loading).toBe(false);
    });

    it('has error null by default', () => {
      const store = useAuthStore();
      expect(store.error).toBeNull();
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Getters
  // ────────────────────────────────────────────────────────────────────────

  describe('getters', () => {
    describe('isAuthenticated', () => {
      it('returns false when accessToken is empty', () => {
        const store = useAuthStore();
        store.accessToken = '';
        store.user = { id: 'u1' };

        expect(store.isAuthenticated).toBe(false);
      });

      it('returns false when user is null', () => {
        const store = useAuthStore();
        store.accessToken = 'some-token';
        store.user = null;

        expect(store.isAuthenticated).toBe(false);
      });

      it('returns false when both are empty/null', () => {
        const store = useAuthStore();

        expect(store.isAuthenticated).toBe(false);
      });

      it('returns true when both accessToken and user are set', () => {
        const store = useAuthStore();
        store.accessToken = 'valid-token';
        store.user = { id: 'u1', email: 'test@test.com' };

        expect(store.isAuthenticated).toBe(true);
      });
    });
  });

  // ────────────────────────────────────────────────────────────────────────
  // Actions
  // ────────────────────────────────────────────────────────────────────────

  describe('actions', () => {
    describe('login', () => {
      it('sets tokens and user on successful login', async () => {
        const { apis } = await import('@src/services/mapex');
        const mockUser = { id: 'u1', email: 'admin@mapex.global', firstName: 'Admin' };
        const mockResponse = {
          access_token: 'access-123',
          refresh_token: 'refresh-456',
          user: mockUser,
        };

        (apis.mapexOS as any).auth = {
          login: vi.fn().mockResolvedValue(mockResponse),
        };

        const store = useAuthStore();
        await store.login('admin@mapex.global', 'mapex123', true);

        expect(store.accessToken).toBe('access-123');
        expect(store.refreshToken).toBe('refresh-456');
        expect(store.user).toEqual(mockUser);
        expect(store.keepConnected).toBe(true);
      });

      it('calls API with correct payload', async () => {
        const { apis } = await import('@src/services/mapex');
        const mockLogin = vi.fn().mockResolvedValue({
          access_token: 'tok',
          refresh_token: 'ref',
          user: { id: 'u1' },
        });
        (apis.mapexOS as any).auth = { login: mockLogin };

        const store = useAuthStore();
        await store.login('user@test.com', 'pass123', false);

        expect(mockLogin).toHaveBeenCalledWith({
          email: 'user@test.com',
          password: 'pass123',
          keepConnected: false,
        });
      });

      it('persists tokens after successful login', async () => {
        const { apis } = await import('@src/services/mapex');
        const storage = (await import('@utils/storages')).default;

        (apis.mapexOS as any).auth = {
          login: vi.fn().mockResolvedValue({
            access_token: 'access-123',
            refresh_token: 'refresh-456',
            user: { id: 'u1' },
          }),
        };

        const store = useAuthStore();
        await store.login('admin@mapex.global', 'mapex123', true);

        expect(storage.local.set).toHaveBeenCalledWith('auth_tokens', {
          accessToken: 'access-123',
          refreshToken: 'refresh-456',
          user: { id: 'u1' },
          keepConnected: true,
        });
      });

      it('uses session storage when keepConnected is false', async () => {
        const { apis } = await import('@src/services/mapex');
        const storage = (await import('@utils/storages')).default;

        (apis.mapexOS as any).auth = {
          login: vi.fn().mockResolvedValue({
            access_token: 'access-123',
            refresh_token: 'refresh-456',
            user: { id: 'u1' },
          }),
        };

        const store = useAuthStore();
        await store.login('admin@mapex.global', 'mapex123', false);

        expect(storage.session.set).toHaveBeenCalledWith('auth_tokens', {
          accessToken: 'access-123',
          refreshToken: 'refresh-456',
          user: { id: 'u1' },
          keepConnected: false,
        });
      });

      it('throws and logs error when API fails', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          login: vi.fn().mockRejectedValue(new Error('Invalid credentials')),
        };

        const store = useAuthStore();

        await expect(store.login('bad@email.com', 'wrong', true))
          .rejects.toThrow('Invalid credentials');
      });

      it('does not set tokens when API fails', async () => {
        const { apis } = await import('@src/services/mapex');

        (apis.mapexOS as any).auth = {
          login: vi.fn().mockRejectedValue(new Error('Fail')),
        };

        const store = useAuthStore();

        try {
          await store.login('bad@email.com', 'wrong', true);
        } catch {
          // expected
        }

        expect(store.accessToken).toBe('');
        expect(store.user).toBeNull();
      });
    });

    describe('logout', () => {
      it('clears store state', () => {
        const store = useAuthStore();
        store.accessToken = 'some-token';
        store.refreshToken = 'some-refresh';
        store.user = { id: 'u1' };
        store.keepConnected = true;

        store.logout();

        expect(store.accessToken).toBe('');
        expect(store.refreshToken).toBe('');
        expect(store.user).toBeNull();
        expect(store.keepConnected).toBe(false);
      });

      it('removes tokens from both storages', async () => {
        const storage = (await import('@utils/storages')).default;

        const store = useAuthStore();
        store.accessToken = 'tok';
        store.logout();

        expect(storage.local.remove).toHaveBeenCalledWith('auth_tokens');
        expect(storage.session.remove).toHaveBeenCalledWith('auth_tokens');
      });

      it('clears organization store coverage', () => {
        const store = useAuthStore();
        store.logout();

        expect(mockOrgStore.clearCoverage).toHaveBeenCalled();
      });

      it('clears permission store', () => {
        const store = useAuthStore();
        store.logout();

        expect(mockPermStore.clearPermissions).toHaveBeenCalled();
      });
    });

    describe('persistTokens', () => {
      it('stores to localStorage when keepConnected is true', async () => {
        const storage = (await import('@utils/storages')).default;

        const store = useAuthStore();
        store.accessToken = 'access-tok';
        store.refreshToken = 'refresh-tok';
        store.user = { id: 'u1' };
        store.keepConnected = true;

        store.persistTokens();

        expect(storage.local.set).toHaveBeenCalledWith('auth_tokens', {
          accessToken: 'access-tok',
          refreshToken: 'refresh-tok',
          user: { id: 'u1' },
          keepConnected: true,
        });
      });

      it('stores to sessionStorage when keepConnected is false', async () => {
        const storage = (await import('@utils/storages')).default;

        const store = useAuthStore();
        store.accessToken = 'access-tok';
        store.refreshToken = 'refresh-tok';
        store.user = { id: 'u1' };
        store.keepConnected = false;

        store.persistTokens();

        expect(storage.session.set).toHaveBeenCalledWith('auth_tokens', {
          accessToken: 'access-tok',
          refreshToken: 'refresh-tok',
          user: { id: 'u1' },
          keepConnected: false,
        });
      });
    });

    describe('hydrateFromStorage', () => {
      it('restores from localStorage when data exists', async () => {
        const storage = (await import('@utils/storages')).default;
        const storedData = {
          accessToken: 'stored-access',
          refreshToken: 'stored-refresh',
          user: { id: 'u1', email: 'test@test.com' },
          keepConnected: true,
        };

        vi.mocked(storage.local.get).mockReturnValue(storedData);

        const store = useAuthStore();
        const result = store.hydrateFromStorage();

        expect(result).toBe(true);
        expect(store.accessToken).toBe('stored-access');
        expect(store.refreshToken).toBe('stored-refresh');
        expect(store.user).toEqual({ id: 'u1', email: 'test@test.com' });
        expect(store.keepConnected).toBe(true);
      });

      it('falls back to sessionStorage when localStorage is empty', async () => {
        const storage = (await import('@utils/storages')).default;
        const storedData = {
          accessToken: 'session-access',
          refreshToken: 'session-refresh',
          user: { id: 'u2' },
          keepConnected: false,
        };

        vi.mocked(storage.local.get).mockReturnValue(null);
        vi.mocked(storage.session.get).mockReturnValue(storedData);

        const store = useAuthStore();
        const result = store.hydrateFromStorage();

        expect(result).toBe(true);
        expect(store.accessToken).toBe('session-access');
        expect(store.keepConnected).toBe(false);
      });

      it('returns false when no stored data exists', async () => {
        const storage = (await import('@utils/storages')).default;

        vi.mocked(storage.local.get).mockReturnValue(null);
        vi.mocked(storage.session.get).mockReturnValue(null);

        const store = useAuthStore();
        const result = store.hydrateFromStorage();

        expect(result).toBe(false);
        expect(store.accessToken).toBe('');
        expect(store.user).toBeNull();
      });

      it('returns false when stored data has no accessToken', async () => {
        const storage = (await import('@utils/storages')).default;

        vi.mocked(storage.local.get).mockReturnValue({ user: { id: 'u1' } });

        const store = useAuthStore();
        const result = store.hydrateFromStorage();

        expect(result).toBe(false);
      });

      it('returns false when stored data has no user', async () => {
        const storage = (await import('@utils/storages')).default;

        vi.mocked(storage.local.get).mockReturnValue({ accessToken: 'tok' });

        const store = useAuthStore();
        const result = store.hydrateFromStorage();

        expect(result).toBe(false);
      });

      it('defaults keepConnected to false when not stored', async () => {
        const storage = (await import('@utils/storages')).default;

        vi.mocked(storage.local.get).mockReturnValue({
          accessToken: 'tok',
          refreshToken: 'ref',
          user: { id: 'u1' },
        });

        const store = useAuthStore();
        store.hydrateFromStorage();

        expect(store.keepConnected).toBe(false);
      });
    });

    describe('updateTokens', () => {
      it('updates access and refresh tokens', () => {
        const store = useAuthStore();
        store.keepConnected = true;

        store.updateTokens('new-access', 'new-refresh');

        expect(store.accessToken).toBe('new-access');
        expect(store.refreshToken).toBe('new-refresh');
      });

      it('re-persists tokens after update', async () => {
        const storage = (await import('@utils/storages')).default;

        const store = useAuthStore();
        store.user = { id: 'u1' };
        store.keepConnected = true;

        store.updateTokens('new-access', 'new-refresh');

        expect(storage.local.set).toHaveBeenCalledWith('auth_tokens', {
          accessToken: 'new-access',
          refreshToken: 'new-refresh',
          user: { id: 'u1' },
          keepConnected: true,
        });
      });
    });
  });
});
