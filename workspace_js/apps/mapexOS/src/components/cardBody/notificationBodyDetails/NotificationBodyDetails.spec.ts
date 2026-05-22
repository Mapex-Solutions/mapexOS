import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import NotificationBodyDetails from './NotificationBodyDetails.vue';

function factory(notification: any) {
  return mountWithPlugins(NotificationBodyDetails, {
    props: { notification },
  });
}

describe('NotificationBodyDetails', () => {
  it('renders without errors for slack type', () => {
    const wrapper = factory({
      type: 'slack',
      data: { workspace: 'Acme Corp', channelsName: ['#alerts', '#general'] },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('renders without errors for email type', () => {
    const wrapper = factory({
      type: 'email',
      data: { from: 'noreply@test.com', to: ['a@test.com', 'b@test.com'] },
    });
    expect(wrapper.exists()).toBe(true);
  });

  it('getMainIcon returns correct icon for slack', () => {
    const wrapper = factory({
      type: 'slack',
      data: { workspace: 'W', channelsName: [] },
    });
    expect(wrapper.vm.getMainIcon('slack')).toBe('mdi-domain');
  });

  it('getMainIcon returns correct icon for teams', () => {
    const wrapper = factory({
      type: 'teams',
      data: { teamName: 'T', channelsName: [] },
    });
    expect(wrapper.vm.getMainIcon('teams')).toBe('mdi-account-group');
  });

  it('getMainColor returns purple-6 for slack', () => {
    const wrapper = factory({
      type: 'slack',
      data: { workspace: 'W', channelsName: [] },
    });
    expect(wrapper.vm.getMainColor('slack')).toBe('purple-6');
  });

  it('getMainLabel returns WORKSPACE for slack', () => {
    const wrapper = factory({
      type: 'slack',
      data: { workspace: 'W', channelsName: [] },
    });
    expect(wrapper.vm.getMainLabel('slack')).toBe('WORKSPACE');
  });

  it('getMainValue returns workspace name for slack', () => {
    const wrapper = factory({
      type: 'slack',
      data: { workspace: 'Acme Corp', channelsName: [] },
    });
    expect(wrapper.vm.getMainValue(wrapper.props('notification'))).toBe('Acme Corp');
  });

  it('visibleEmails returns first 4 emails by default', () => {
    const emails = ['a@t.com', 'b@t.com', 'c@t.com', 'd@t.com', 'e@t.com'];
    const wrapper = factory({
      type: 'email',
      data: { from: 'x@t.com', to: emails },
    });
    expect(wrapper.vm.visibleEmails).toHaveLength(4);
  });

  it('hasMoreEmails is true when more than 4 emails', () => {
    const emails = ['a@t.com', 'b@t.com', 'c@t.com', 'd@t.com', 'e@t.com'];
    const wrapper = factory({
      type: 'email',
      data: { from: 'x@t.com', to: emails },
    });
    expect(wrapper.vm.hasMoreEmails).toBe(true);
  });

  it('hasMoreEmails is false when 4 or fewer emails', () => {
    const wrapper = factory({
      type: 'email',
      data: { from: 'x@t.com', to: ['a@t.com', 'b@t.com'] },
    });
    expect(wrapper.vm.hasMoreEmails).toBe(false);
  });

  it('allItems returns channelsName for slack', () => {
    const wrapper = factory({
      type: 'slack',
      data: { workspace: 'W', channelsName: ['#a', '#b'] },
    });
    expect(wrapper.vm.allItems).toEqual(['#a', '#b']);
  });

  it('allItems returns chatNames for telegram', () => {
    const wrapper = factory({
      type: 'telegram',
      data: { botName: 'bot', chatNames: ['chat1', 'chat2'] },
    });
    expect(wrapper.vm.allItems).toEqual(['chat1', 'chat2']);
  });

  it('getChipColor returns purple-2 for slack', () => {
    const wrapper = factory({
      type: 'slack',
      data: { workspace: 'W', channelsName: [] },
    });
    expect(wrapper.vm.getChipColor('slack')).toBe('purple-2');
  });

  it('getConfigTitle returns correct title for each type', () => {
    const wrapper = factory({
      type: 'push',
      data: { appName: 'App', deviceCount: 100 },
    });
    expect(wrapper.vm.getConfigTitle('push')).toBe('NOTIFICA\u00C7\u00D5ES PUSH');
    expect(wrapper.vm.getConfigTitle('slack')).toBe('CONFIGURA\u00C7\u00C3O');
    expect(wrapper.vm.getConfigTitle('teams')).toBe('INTEGRA\u00C7\u00C3O');
  });
});
