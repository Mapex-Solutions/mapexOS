/**
 * Get valus nested on object / array
 * @param obj: Record<any, any> - Object data
 * @param path: string - Full path to value
 * @param defaultValue - If not exist the property set the value
 */
export async function getByPathAsync(obj: Record<any, any>, path: string, defaultValue?: any) {
  const undefinedValue = undefined;
  const keys = path.split('.');

  if (!keys.every((el) => String(el).length)) return obj;

  let result = obj;
  let index = 0;

  for (; index < keys.length;) {
    const key = keys[index];

    if (
      Array.isArray(result) ||
      (typeof result === 'object' && key in result)
    ) {
      // Is object
      if (!Array.isArray(result)) {
        result = result[key];
        keys.shift();
        index = -1;
      }

      // Is array
      if (Array.isArray(result)) {
        const items = [];

        for await (const item of result) {
          const value = await getByPath(item, keys.join('.'), undefinedValue);
          if (value !== undefined) items.push(value);
        }

        result = items.length ? items.flat() : defaultValue;
        break;
      }
    } else {
      result = defaultValue;
    }

    index += 1;
  }

  return result;
}


/**
 * Get valus nested on object / array
 * @param obj: Record<any, any> - Object data
 * @param path: string - Full path to value
 * @param defaultValue - If not exist the property set the value
 */
export function getByPath(obj: Record<any, any>, path: string, defaultValue?: any) {

  const undefinedValue = undefined;
  const keys = path.split('.');

  if (!keys.every((el) => String(el).length)) return obj;

  let result = obj;
  let index = 0;

  for (; index < keys.length;) {
    const key = keys[index];

    if (
      Array.isArray(result) ||
      (typeof result === 'object' && key in result)
    ) {
      // Is object
      if (!Array.isArray(result)) {
        result = result[key];
        keys.shift();
        index = -1;
      }

      // Is array
      if (Array.isArray(result)) {
        const items = [];

        for (const item of result) {
          const value = getByPath(item, keys.join('.'), undefinedValue);
          if (value !== undefined) items.push(value);
        }

        result = items.length ? items.flat() : defaultValue;
        break;
      }
    } else {
      result = defaultValue;
    }

    index += 1;
  }

  return result;
}
