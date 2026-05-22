/**
 * Test helpers for mounting Vue components with all required plugins.
 */
import { mount, shallowMount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import { vi } from 'vitest';
import type { Component } from 'vue';

/**
 * Options for mountWithPlugins
 */
interface MountOptions {
   
  props?: Record<string, any>;
  slots?: Record<string, string | Component>;
  shallow?: boolean;
  piniaState?: Record<string, unknown>;
  stubs?: Record<string, boolean | Component>;
  attachTo?: HTMLElement | string;
}

/**
 * Mount a component with Pinia + common stubs pre-configured.
 * Uses shallowMount by default for faster tests.
 *
 * @param {Component} component - Vue component to mount
 * @param {MountOptions} options - Mount options
 * @returns VueWrapper
 */
export function mountWithPlugins(component: Component, options: MountOptions = {}) {
  const { shallow = true, stubs } = options;

  const mountFn = shallow ? shallowMount : mount;

  return mountFn(component, {
    ...(options.props != null && { props: options.props }),
    ...(options.slots != null && { slots: options.slots }),
    ...(options.attachTo != null && { attachTo: options.attachTo }),
    global: {
      plugins: [
        createTestingPinia({
          createSpy: vi.fn,
          ...(options.piniaState != null && { initialState: options.piniaState }),
        }),
      ],
      stubs: {
        ...stubs,
      },
    },
  });
}

/**
 * Create a mock for any composable that returns computed refs.
 * Returns an object where each value is a ref-like { value: X }.
 *
 * @param {Record<string, unknown>} overrides - Values to return
 * @returns Mocked composable return value
 */
export function createMockTranslations(overrides: Record<string, unknown> = {}) {
  const handler: ProxyHandler<Record<string, unknown>> = {
    get(_target, prop) {
      if (prop === 'value') return prop;
      if (typeof prop === 'string' && prop in overrides) {
        return overrides[prop];
      }
      // Return a nested proxy for deep access (t.page.title.value → 'title')
      return new Proxy({ value: String(prop) }, handler);
    },
  };

  return new Proxy({}, handler);
}
