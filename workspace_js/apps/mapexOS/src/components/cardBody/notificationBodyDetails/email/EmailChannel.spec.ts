import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import EmailChannel from './EmailChannel.vue';

function factory(channel: { type: 'email'; from: string; to: string[] }) {
  return mountWithPlugins(EmailChannel, {
    props: { channel },
  });
}

describe('EmailChannel', () => {
  it('renders without errors', () => {
    const wrapper = factory({ type: 'email', from: 'sender@test.com', to: ['a@test.com'] });
    expect(wrapper.exists()).toBe(true);
  });

  it('displayList shows first 2 emails', () => {
    const wrapper = factory({
      type: 'email',
      from: 'sender@test.com',
      to: ['a@test.com', 'b@test.com', 'c@test.com', 'd@test.com'],
    });
    expect(wrapper.vm.displayList).toEqual(['a@test.com', 'b@test.com']);
  });

  it('extraList contains emails beyond first 2', () => {
    const wrapper = factory({
      type: 'email',
      from: 'sender@test.com',
      to: ['a@test.com', 'b@test.com', 'c@test.com'],
    });
    expect(wrapper.vm.extraList).toEqual(['c@test.com']);
  });

  it('extraCount is 0 when 2 or fewer emails', () => {
    const wrapper = factory({
      type: 'email',
      from: 'sender@test.com',
      to: ['a@test.com', 'b@test.com'],
    });
    expect(wrapper.vm.extraCount).toBe(0);
  });

  it('extraCount reflects number of hidden emails', () => {
    const wrapper = factory({
      type: 'email',
      from: 'sender@test.com',
      to: ['a@test.com', 'b@test.com', 'c@test.com', 'd@test.com'],
    });
    expect(wrapper.vm.extraCount).toBe(2);
  });
});
