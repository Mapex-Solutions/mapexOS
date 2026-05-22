package middlewares

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/auth"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/requestValidation"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"

	"http_gateway/src/bootstrap"
	dsDto "http_gateway/src/modules/datasources/application/dtos"
	dsPort "http_gateway/src/modules/datasources/application/ports"

	"http_gateway/src/modules/events/application/dtos"
	eventsPort "http_gateway/src/modules/events/application/ports"
)

// CustomAuthMiddleware authorizes the request based on the DataSource auth settings.
//
// Following Hexagonal Architecture, this middleware accepts the DataSourceServicePort
// interface rather than a concrete service implementation.
//
// It expects an EvenIdentificationDto in the query (e.g., ds parameter) and applies
// authentication based on the data source configuration (OAuth2, JWT, API Key, IP whitelist, or none).
//
// When authentication fails, it publishes a raw event with success=false to events.raw
// for security monitoring purposes.
//
// Parameters:
//   - dsService: Data source service port interface for retrieving data source configurations
//   - eventService: Event service port interface for publishing auth failure events
//   - m: Service-specific metrics for auth instrumentation
//
// Returns:
//   - A Fiber handler function that validates authentication before allowing access
func CustomAuthMiddleware(dsService dsPort.DataSourceServicePort, eventService eventsPort.EventServicePort, m *bootstrap.HttpGatewayMetrics) fiber.Handler {

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// The data already validated by the validation middleware
		dataSource, err := getDataSource(c, dsService)

		if err != nil {
			return err
		}

		// Runtime gate (TKT-2026-0036): disabled DataSources reject all traffic
		// before any auth cost. Covers both /events and /heartbeat (single change
		// point). The "disabled" metric label is added to the existing 3 auth
		// metrics — no new metric definitions.
		// Enabled is *bool: nil or false → disabled (defensive default).
		if dataSource.Enabled == nil || !*dataSource.Enabled {
			errMsg := "DataSource is disabled"
			eventService.PublishAuthFailure(dataSource, parseEventBody(c), uuid.New().String(), errMsg)
			duration := time.Since(start).Seconds()
			m.EventAuthTotal.WithLabelValues("disabled", "failure").Inc()
			m.EventAuthDuration.WithLabelValues("disabled").Observe(duration)
			m.EventAuthFailures.WithLabelValues("disabled").Inc()
			return &customErrors.ServerCustomError{Code: status.FORBIDDEN, Errors: []string{errMsg}}
		}

		authType := dataSource.Auth.Type

		switch authType {

		case "oauth2":
			err = checkJWKS(c, dataSource, eventService)

		case "jwt":
			err = checkJWT(c, dataSource, eventService)

		case "apiKey":
			err = checkApiKey(c, dataSource, eventService)

		case "ip_whitelist":
			err = checkIPWhiteList(c, dataSource, eventService)

		case "none":
			err = none(c, dataSource)

		default:
			errMsg := "Auth type not supported"
			eventService.PublishAuthFailure(dataSource, parseEventBody(c), uuid.New().String(), errMsg)
			duration := time.Since(start).Seconds()

			// Metric updates for unsupported auth type
			m.EventAuthTotal.WithLabelValues(authType, "failure").Inc()
			m.EventAuthDuration.WithLabelValues(authType).Observe(duration)
			m.EventAuthFailures.WithLabelValues(authType).Inc()

			return &customErrors.ServerCustomError{Code: status.UNAUTHORIZED, Errors: []string{errMsg}}
		}

		duration := time.Since(start).Seconds()

		if err != nil {
			// Metrics: auth validation failed — count attempt as failure, record latency, track security event
			m.EventAuthTotal.WithLabelValues(authType, "failure").Inc()
			m.EventAuthDuration.WithLabelValues(authType).Observe(duration)
			m.EventAuthFailures.WithLabelValues(authType).Inc()
			return err
		}

		// Metrics: auth validation succeeded — count attempt as success, record strategy latency
		m.EventAuthTotal.WithLabelValues(authType, "success").Inc()
		m.EventAuthDuration.WithLabelValues(authType).Observe(duration)
		return nil
	}
}

// parseEventBody parses the request body for inclusion in auth failure events.
// Called lazily only when auth fails to avoid unnecessary parsing on the happy path.
func parseEventBody(c *fiber.Ctx) map[string]any {
	var event map[string]any
	_ = c.BodyParser(&event)
	return event
}

func getDataSource(c *fiber.Ctx, dsService dsPort.DataSourceServicePort) (*dsDto.DataSourceResponse, error) {
	// retrieve the timeout‐aware Context you set in ContextInjector
	ctx := c.UserContext()

	// The data already validated by the validation middleware
	queryData, _ := requestValidation.GetDTO[*dtos.EvenIdentificationDto](c, "queryDTO")

	dataSource, err := dsService.GetDataSourceById(ctx, queryData.Ds)
	if err != nil {
		return nil, err
	}

	return dataSource, nil
}

func none(c *fiber.Ctx, dataSource *dsDto.DataSourceResponse) error {
	c.Locals("dataSource", dataSource)
	return c.Next()
}

func checkIPWhiteList(c *fiber.Ctx, dataSource *dsDto.DataSourceResponse, eventService eventsPort.EventServicePort) error {
	canAccess := auth.ValidateIPWhitelist(c, dataSource.Auth.IPWhitelist.CIDRs)

	if !canAccess {
		errMsg := "Unauthorized - IP Address not allowed"
		eventService.PublishAuthFailure(dataSource, parseEventBody(c), uuid.New().String(), errMsg)
		return &customErrors.ServerCustomError{Code: status.UNAUTHORIZED, Errors: []string{errMsg}}
	}

	c.Locals("dataSource", dataSource)
	return c.Next()
}

// checkJWKS validates the JWT token provided in the Authorization header of the HTTP request
// using the JWKS URL specified in the DataSource. It extracts the token from the header,
// verifies its validity and expiration, and returns an error if the token is missing, invalid, or expired.
//
// Parameters:
//   - c: Fiber context for the current HTTP request.
//   - dataSource: DataSourceResponse containing the OAuth2 auth settings with the JWKS URL.
//   - event: The event payload to include in failure events.
//   - eventTrackerId: UUID for tracking event across the pipeline.
//   - eventService: Event service for publishing auth failures.
//
// Returns:
//   - nil if the token is valid and not expired.
//   - customErrors.ServerCustomError with status.UNAUTHORIZED if the token is missing, invalid, or expired.
func checkJWKS(c *fiber.Ctx, dataSource *dsDto.DataSourceResponse, eventService eventsPort.EventServicePort) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		errMsg := "missing Authorization header"
		eventService.PublishAuthFailure(dataSource, parseEventBody(c), uuid.New().String(), errMsg)
		return &customErrors.ServerCustomError{Code: status.UNAUTHORIZED, Errors: []string{errMsg}}
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	/** Check if the token is valid and not expired */
	if _, _, err := auth.ParseJWTTokenWithJWKS(
		tokenString,
		dataSource.Auth.OAuth2.JWKSURL,
	); err != nil {
		errMsg := "invalid token: " + err.Error()
		eventService.PublishAuthFailure(dataSource, parseEventBody(c), uuid.New().String(), errMsg)
		return &customErrors.ServerCustomError{Code: status.UNAUTHORIZED, Errors: []string{"invalid token", err.Error()}}
	}

	c.Locals("dataSource", dataSource)
	return c.Next()
}

// checkJWT validates the JWT token provided in the configured header of the HTTP request.
// It extracts the token from the header, verifies its validity and expiration, and returns an error if the token is missing, invalid, or expired.
//
// Parameters:
//   - c: Fiber context for the current HTTP request.
//   - dataSource: DataSourceResponse containing the JWT auth settings.
//   - event: The event payload to include in failure events.
//   - eventTrackerId: UUID for tracking event across the pipeline.
//   - eventService: Event service for publishing auth failures.
//
// Returns:
//   - nil if the token is valid and not expired.
//   - customErrors.ServerCustomError with status.UNAUTHORIZED if the token is missing, invalid, or expired.
func checkJWT(c *fiber.Ctx, dataSource *dsDto.DataSourceResponse, eventService eventsPort.EventServicePort) error {

	// Determine header name - use configured or default to "Authorization"
	headerName := "Authorization"
	if dataSource.Auth.JWT.HeaderName != nil && *dataSource.Auth.JWT.HeaderName != "" {
		headerName = *dataSource.Auth.JWT.HeaderName
	}

	/** Get the JWT token from the configured header */
	authHeader := c.Get(headerName)
	if authHeader == "" {
		errMsg := "missing " + headerName + " header"
		eventService.PublishAuthFailure(dataSource, parseEventBody(c), uuid.New().String(), errMsg)
		return &customErrors.ServerCustomError{Code: status.UNAUTHORIZED, Errors: []string{errMsg}}
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	/** Check if the token is valid and not expired */
	if _, _, err := auth.ParseJWTTokenWithSecret(
		tokenString,
		dataSource.Auth.JWT.Secret,
		dataSource.Auth.JWT.Algorithms,
	); err != nil {
		errMsg := "invalid token: " + err.Error()
		eventService.PublishAuthFailure(dataSource, parseEventBody(c), uuid.New().String(), errMsg)
		return &customErrors.ServerCustomError{Code: status.UNAUTHORIZED, Errors: []string{"invalid token", err.Error()}}
	}

	c.Locals("dataSource", dataSource)
	return c.Next()
}

// checkApiKey validates access to the current request using the API-Key
// configuration found in the provided DataSource.
//
// It delegates the credential extraction and comparison to auth.ValidateAPIKey,
// which reads the candidate key from the request according to
// dataSource.Auth.APIKey.Type (e.g., "header" or "query") and
// dataSource.Auth.APIKey.FieldName, and compares it with
// dataSource.Auth.APIKey.Key.
//
// Behavior:
//   - On success: stores the DataSource in the request context via
//     c.Locals("dataSource", dataSource) and calls the next handler
//     (c.Next()).
//   - On failure: publishes auth failure event and returns a
//     customErrors.ServerCustomError with status.UNAUTHORIZED.
//
// Params:
//   - c: Fiber context for the current HTTP request.
//   - dataSource: DataSourceResponse containing the API-Key auth settings.
//   - event: The event payload to include in failure events.
//   - eventTrackerId: UUID for tracking event across the pipeline.
//   - eventService: Event service for publishing auth failures.
//
// Returns:
//
//   - error from the next handler on success, or an Unauthorized error if the
//     API key is missing/invalid.
func checkApiKey(c *fiber.Ctx, dataSource *dsDto.DataSourceResponse, eventService eventsPort.EventServicePort) error {
	canAccess := auth.ValidateAPIKey(
		c,
		dataSource.Auth.APIKey.Key,
		dataSource.Auth.APIKey.Type,
		dataSource.Auth.APIKey.FieldName,
	)

	if !canAccess {
		errMsg := "Unauthorized - Invalid API Key"
		eventService.PublishAuthFailure(dataSource, parseEventBody(c), uuid.New().String(), errMsg)
		return &customErrors.ServerCustomError{Code: status.UNAUTHORIZED, Errors: []string{errMsg}}
	}

	c.Locals("dataSource", dataSource)
	return c.Next()
}
