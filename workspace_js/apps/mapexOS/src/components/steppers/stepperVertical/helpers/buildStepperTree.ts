import type { StepperVerticalItem, GroupedStep } from '../interfaces';

/**
 * Collapses a flat list of leaves into the vertical stepper's grouped tree:
 * consecutive leaves sharing a group id nest under one group header; ungrouped
 * leaves stay flat. Leaf order is preserved.
 *
 * Navigation and the current-step index still operate on the flat leaf list —
 * this only shapes what StepperVertical renders. Passing leaves with no `group`
 * yields the flat list unchanged (the default, backward-compatible behavior).
 * @param {GroupedStep[]} leaves - Flat, ordered steps, each optionally grouped.
 * @returns {StepperVerticalItem[]} The nested tree for StepperVertical's `steps`.
 */
export function buildStepperTree(leaves: GroupedStep[]): StepperVerticalItem[] {
  const tree: StepperVerticalItem[] = [];
  let currentGroupId: string | null = null;

  for (const leaf of leaves) {
    const item: StepperVerticalItem = { title: leaf.title, description: leaf.description, icon: leaf.icon };

    if (leaf.group) {
      const parent = tree[tree.length - 1];
      if (currentGroupId === leaf.group.id && parent) {
        parent.children = parent.children ? [...parent.children, item] : [item];
      } else {
        tree.push({ title: leaf.group.label, icon: leaf.group.icon, description: '', children: [item] });
        currentGroupId = leaf.group.id;
      }
    } else {
      tree.push(item);
      currentGroupId = null;
    }
  }

  return tree;
}
