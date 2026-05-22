package constants

import "time"

/**
 * Cache configuration for DataSource
 */

// DataSourceCacheTTL defines the time-to-live for DataSource cache entries.
// 24 hours is appropriate since datasources change very infrequently.
const DataSourceCacheTTL = 24 * time.Hour
