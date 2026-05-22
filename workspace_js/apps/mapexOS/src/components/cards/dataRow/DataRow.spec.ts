import { describe, it, expect } from 'vitest';
import { mountWithPlugins } from '@src/test/helpers';
import DataRow from './DataRow.vue';
import type { DataRowColumn, DataRowActionConfig } from './interfaces';

const baseColumns: DataRowColumn[] = [
  { key: 'avatar', label: '', type: 'avatar', visible: 'always', width: 60 },
  { key: 'name', label: 'Name', type: 'text', visible: 'always' },
  { key: 'status', label: 'Status', type: 'badge', visible: 'always' },
  { key: 'protocol.type', label: 'Protocol', type: 'chip', visible: 'laptop' },
  { key: 'uuid', label: 'UUID', type: 'code', visible: 'desktop' },
];

const baseData = {
  id: '1',
  avatar: 'A',
  name: 'Test Asset',
  status: 'active',
  protocol: { type: 'mqtt' },
  uuid: 'abc-123',
};

function factory(overrides: Record<string, unknown> = {}) {
  return mountWithPlugins(DataRow, {
    props: {
      data: baseData,
      columns: baseColumns,
      ...overrides,
    },
  });
}

describe('DataRow', () => {
  it('renders without errors', () => {
    const wrapper = factory();
    expect(wrapper.exists()).toBe(true);
  });

  it('tableColumns computed maps columns to q-table format', () => {
    const wrapper = factory();
    const tc = wrapper.vm.tableColumns;

    expect(tc).toHaveLength(baseColumns.length);
    expect(tc[0].name).toBe('avatar');
    expect(tc[0].field).toBe('avatar');
    expect(tc[0].align).toBe('left');
  });

  it('tableColumns includes width style when column has width', () => {
    const wrapper = factory();
    const avatarCol = wrapper.vm.tableColumns.find((c: any) => c.name === 'avatar');
    expect(avatarCol.style).toContain('width: 60px');
  });

  it('showEdit defaults to true when actions prop is undefined', () => {
    const wrapper = factory();
    expect(wrapper.vm.showEdit).toBe(true);
  });

  it('showView defaults to true when actions prop is undefined', () => {
    const wrapper = factory();
    expect(wrapper.vm.showView).toBe(true);
  });

  it('showDelete defaults to true when actions prop is undefined', () => {
    const wrapper = factory();
    expect(wrapper.vm.showDelete).toBe(true);
  });

  it('showEdit respects actions config set to false', () => {
    const actions: DataRowActionConfig = { showEdit: false };
    const wrapper = factory({ actions });
    expect(wrapper.vm.showEdit).toBe(false);
  });

  it('showView respects actions config set to false', () => {
    const actions: DataRowActionConfig = { showView: false };
    const wrapper = factory({ actions });
    expect(wrapper.vm.showView).toBe(false);
  });

  it('showDelete respects actions config set to false', () => {
    const actions: DataRowActionConfig = { showDelete: false };
    const wrapper = factory({ actions });
    expect(wrapper.vm.showDelete).toBe(false);
  });

  it('visibleCustomActions filters by condition', () => {
    const actions: DataRowActionConfig = {
      customActions: [
        { key: 'a', label: 'A', icon: 'star', condition: () => true },
        { key: 'b', label: 'B', icon: 'star', condition: () => false },
        { key: 'c', label: 'C', icon: 'star' },
      ],
    };
    const wrapper = factory({ actions });
    const visible = wrapper.vm.visibleCustomActions;
    expect(visible).toHaveLength(2);
    expect(visible.map((a: any) => a.key)).toEqual(['a', 'c']);
  });

  it('avatarColumn finds column with type avatar', () => {
    const wrapper = factory();
    expect(wrapper.vm.avatarColumn?.key).toBe('avatar');
  });

  it('nameColumn finds column with key name', () => {
    const wrapper = factory();
    expect(wrapper.vm.nameColumn?.key).toBe('name');
  });

  it('statusColumn finds column with key status', () => {
    const wrapper = factory();
    expect(wrapper.vm.statusColumn?.key).toBe('status');
  });

  it('mobileExpandableColumns excludes avatar, name, and status', () => {
    const wrapper = factory();
    const expandable = wrapper.vm.mobileExpandableColumns;
    const keys = expandable.map((c: any) => c.key);
    expect(keys).not.toContain('avatar');
    expect(keys).not.toContain('name');
    expect(keys).not.toContain('status');
    expect(keys).toContain('protocol.type');
    expect(keys).toContain('uuid');
  });

  it('getColumnValue resolves nested dot notation keys', () => {
    const wrapper = factory();
    const col = { key: 'protocol.type' } as DataRowColumn;
    expect(wrapper.vm.getColumnValue(col)).toBe('mqtt');
  });

  it('getColumnValue resolves simple keys', () => {
    const wrapper = factory();
    const col = { key: 'name' } as DataRowColumn;
    expect(wrapper.vm.getColumnValue(col)).toBe('Test Asset');
  });

  it('toggleMobileExpand toggles mobileExpanded state', () => {
    const wrapper = factory();
    expect(wrapper.vm.mobileExpanded).toBe(false);
    wrapper.vm.toggleMobileExpand();
    expect(wrapper.vm.mobileExpanded).toBe(true);
    wrapper.vm.toggleMobileExpand();
    expect(wrapper.vm.mobileExpanded).toBe(false);
  });

  it('toggleMobileExpand does nothing when expandOnClick is false', () => {
    const wrapper = factory({ expandOnClick: false });
    wrapper.vm.toggleMobileExpand();
    expect(wrapper.vm.mobileExpanded).toBe(false);
  });
});
