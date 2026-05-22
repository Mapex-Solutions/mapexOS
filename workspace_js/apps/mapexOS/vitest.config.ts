import { defineConfig } from 'vitest/config';
import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';

export default defineConfig({
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'happy-dom',
    setupFiles: ['./src/test/setup.ts'],
    include: ['src/**/*.spec.ts'],
    coverage: {
      provider: 'v8',
      include: [
        'src/components/**/*.vue',
        'src/stores/**/*.ts',
        'src/utils/**/*.ts',
      ],
    },
  },
  resolve: {
    alias: {
      '@src': resolve(__dirname, 'src'),
      '@components': resolve(__dirname, 'src/components'),
      '@stores': resolve(__dirname, 'src/stores'),
      '@composables': resolve(__dirname, 'src/composables'),
      '@utils': resolve(__dirname, 'src/utils'),
      '@services': resolve(__dirname, 'src/services'),
      '@interfaces': resolve(__dirname, 'src/interfaces'),
      src: resolve(__dirname, 'src'),
      'pages': resolve(__dirname, 'src/pages'),
      'components': resolve(__dirname, 'src/components'),
      'stores': resolve(__dirname, 'src/stores'),
      'boot': resolve(__dirname, 'src/boot'),
      'assets': resolve(__dirname, 'src/assets'),
      'monaco-editor/esm': resolve(__dirname, 'src/test/__mocks__/monaco-editor.ts'),
      'monaco-editor': resolve(__dirname, 'src/test/__mocks__/monaco-editor.ts'),
    },
  },
});
