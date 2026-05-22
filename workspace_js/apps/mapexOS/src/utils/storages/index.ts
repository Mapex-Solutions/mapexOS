import { isEmpty } from 'lodash';

function getLocalStorage(): Storage | null {
  return typeof window !== 'undefined' && window.localStorage
    ? window.localStorage
    : null;
}

function getSessionStorage(): Storage | null {
  return typeof window !== 'undefined' && window.sessionStorage
    ? window.sessionStorage
    : null;
}

export default {
  local: {
    set: (key: string, data: any) => {
      const storage = getLocalStorage();
      if (storage) storage.setItem(key, JSON.stringify(data));
    },
    get: (key: string) => {
      const storage = getLocalStorage();
      if (!storage) return null;
      const value = storage.getItem(key);
      return value ? JSON.parse(value) : null;
    },
    remove: (key: string) => {
      const storage = getLocalStorage();
      if (storage) storage.removeItem(key);
    },
    clear: () => {
      const storage = getLocalStorage();
      if (storage) storage.clear();
    },
    sessionClean: function () {
      const custom = this.get('_custom');
      this.clear();
      if (!isEmpty(custom)) {
        this.set('_custom', custom);
      }
    }
  },
  session: {
    set: (key: string, data: any) => {
      const storage = getSessionStorage();
      if (storage) storage.setItem(key, JSON.stringify(data));
    },
    get: (key: string) => {
      const storage = getSessionStorage();
      if (!storage) return null;
      const value = storage.getItem(key);
      return value ? JSON.parse(value) : null;
    },
    remove: (key: string) => {
      const storage = getSessionStorage();
      if (storage) storage.removeItem(key);
    },
    clear: () => {
      const storage = getSessionStorage();
      if (storage) storage.clear();
    },
    sessionClean: function () {
      this.clear();
    }
  }
};
