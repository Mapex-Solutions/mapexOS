package constants

import "time"

// DefaultBrokerServerTTL applies to broker server certs signed by the
// intermediate CA via POST /sign_server when the request omits ttlDays.
const DefaultBrokerServerTTL = 10 * 365 * 24 * time.Hour
