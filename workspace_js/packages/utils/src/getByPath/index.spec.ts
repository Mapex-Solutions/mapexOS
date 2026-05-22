import { getByPath, getByPathAsync } from './index';

describe('getByPath', () => {
  const testObj = {
    data: {
      event: {
        id: 'event-123',
        name: 'Test Event',
        metadata: {
          type: 'user-action',
          timestamp: 1234567890
        }
      },
      users: [
        { id: 1, name: 'John', active: true },
        { id: 2, name: 'Jane', active: false },
        { id: 3, name: 'Bob', active: true }
      ]
    },
    config: {
      settings: {
        theme: 'dark',
        notifications: true
      }
    }
  };

  describe('basic object navigation', () => {
    it('should get simple nested property', () => {
      expect(getByPath(testObj, 'data.event.id')).toBe('event-123');
    });

    it('should get deeply nested property', () => {
      expect(getByPath(testObj, 'data.event.metadata.type')).toBe('user-action');
    });

    it('should get property from different branch', () => {
      expect(getByPath(testObj, 'config.settings.theme')).toBe('dark');
    });
  });

  describe('array handling', () => {
    it('should extract property from all array items', () => {
      const result = getByPath(testObj, 'data.users.name');
      expect(result).toEqual(['John', 'Jane', 'Bob']);
    });

    it('should extract nested property from array items', () => {
      const result = getByPath(testObj, 'data.users.id');
      expect(result).toEqual([1, 2, 3]);
    });

    it('should handle boolean properties in arrays', () => {
      const result = getByPath(testObj, 'data.users.active');
      expect(result).toEqual([true, false, true]);
    });
  });

  describe('edge cases', () => {
    it('should return undefined for non-existent path', () => {
      expect(getByPath(testObj, 'data.nonexistent.property')).toBeUndefined();
    });

    it('should return default value for non-existent path', () => {
      expect(getByPath(testObj, 'data.nonexistent.property', 'default')).toBe('default');
    });

    it('should return original object for empty path', () => {
      expect(getByPath(testObj, '')).toBe(testObj);
    });

    it('should return original object for invalid path with empty segments', () => {
      expect(getByPath(testObj, 'data..event')).toBe(testObj);
    });
  });

  describe('complex nested arrays', () => {
    const complexObj = {
      departments: [
        {
          name: 'Engineering',
          teams: [
            { name: 'Frontend', members: [{ name: 'Alice' }, { name: 'Bob' }] },
            { name: 'Backend', members: [{ name: 'Charlie' }, { name: 'David' }] }
          ]
        },
        {
          name: 'Design',
          teams: [
            { name: 'UX', members: [{ name: 'Eve' }, { name: 'Frank' }] }
          ]
        }
      ]
    };

    it('should extract from nested arrays', () => {
      const result = getByPath(complexObj, 'departments.teams.name');
      expect(result).toEqual(['Frontend', 'Backend', 'UX']);
    });

    it('should extract from deeply nested arrays', () => {
      const result = getByPath(complexObj, 'departments.teams.members.name');
      expect(result).toEqual(['Alice', 'Bob', 'Charlie', 'David', 'Eve', 'Frank']);
    });
  });

  describe('array with mixed data', () => {
    const mixedObj = {
      items: [
        { type: 'A', value: 1 },
        { type: 'B', value: 2 },
        { type: 'A' }, // missing value
        { type: 'C', value: 3 }
      ]
    };

    it('should only return defined values from arrays', () => {
      const result = getByPath(mixedObj, 'items.value');
      expect(result).toEqual([1, 2, 3]);
    });

    it('should return all values including undefined when explicitly checking', () => {
      const types = getByPath(mixedObj, 'items.type');
      expect(types).toEqual(['A', 'B', 'A', 'C']);
    });
  });
});

describe('getByPathAsync', () => {
  const testObj = {
    data: {
      event: {
        id: 'event-123',
        name: 'Test Event'
      },
      users: [
        { id: 1, name: 'John' },
        { id: 2, name: 'Jane' }
      ]
    }
  };

  describe('basic async functionality', () => {
    it('should get simple nested property asynchronously', async () => {
      const result = await getByPathAsync(testObj, 'data.event.id');
      expect(result).toBe('event-123');
    });

    it('should extract property from array items asynchronously', async () => {
      const result = await getByPathAsync(testObj, 'data.users.name');
      expect(result).toEqual(['John', 'Jane']);
    });

    it('should return default value for non-existent path', async () => {
      const result = await getByPathAsync(testObj, 'data.nonexistent.property', 'default');
      expect(result).toBe('default');
    });

    it('should return undefined for non-existent path without default', async () => {
      const result = await getByPathAsync(testObj, 'data.nonexistent.property');
      expect(result).toBeUndefined();
    });
  });

  describe('async edge cases', () => {
    it('should handle empty path', async () => {
      const result = await getByPathAsync(testObj, '');
      expect(result).toBe(testObj);
    });
  });
});
