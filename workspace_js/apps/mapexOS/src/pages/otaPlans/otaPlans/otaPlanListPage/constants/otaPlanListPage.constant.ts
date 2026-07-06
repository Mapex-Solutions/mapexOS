import type { OtaPlanListFilters } from '../interfaces/otaPlanListPage.interface';

/** Default page size for the plans list */
export const DEFAULT_ITEMS_PER_PAGE = 20;

/** Quick-filter defaults (name search + wire status) */
export const OTA_PLAN_FILTER_DEFAULTS: OtaPlanListFilters = {
	name: undefined,
	status: undefined,
};
