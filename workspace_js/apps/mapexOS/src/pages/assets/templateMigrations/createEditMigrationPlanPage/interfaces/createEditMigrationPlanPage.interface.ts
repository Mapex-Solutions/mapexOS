// CreateEditMigrationPlanPage Interfaces

/**
 * Shared reactive form model for the create/edit migration plan wizard.
 * In edit mode only name, description and scheduleAt are editable; the
 * template/asset fields are immutable after creation.
 */
export interface MigrationPlanForm {
  name: string;
  description: string;
  scheduleAt: string | null;
  fromTemplateId: string | null;
  toTemplateId: string | null;
  assetIds: string[];
}
