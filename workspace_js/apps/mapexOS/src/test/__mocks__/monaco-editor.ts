/**
 * Mock for monaco-editor — prevents worker import errors in Vitest.
 */
export const editor = {
  create: () => ({
    dispose: () => {},
    getValue: () => '',
    setValue: () => {},
    onDidChangeModelContent: () => ({ dispose: () => {} }),
    getModel: () => null,
    updateOptions: () => {},
  }),
  defineTheme: () => {},
  setTheme: () => {},
  createModel: () => ({}),
};

export const languages = {
  register: () => {},
  setMonarchTokensProvider: () => {},
  registerCompletionItemProvider: () => ({ dispose: () => {} }),
};

export const Uri = {
  parse: () => ({}),
};

export const KeyMod = {};
export const KeyCode = {};
export default { editor, languages, Uri, KeyMod, KeyCode };
