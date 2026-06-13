package bootstrap

import (
	"fmt"

	web "github.com/Mapex-Solutions/mapexGoKit/microservices/http/web"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/http/swagger"
)

// InitSwagger sets the OpenAPI document metadata from service configuration.
// Call before modules register their routes so each route records into a
// configured registry.
//
// Document-level fields (description, servers, contact, license, tags) are filled
// here, once per service — route-level details are attached at each route.
func InitSwagger() {
	serviceName, _ := config.GetStringValue("service_name")
	serviceVersion, _ := config.GetStringValue("service_version")
	httpAddress, _ := config.GetStringValue("http_address")
	httpPort, _ := config.GetIntValue("http_port")
	goEnv, _ := config.GetStringValue("go_env")

	swagger.Init(swagger.Config{
		ServiceName: serviceName,
		Version:     serviceVersion,
		Description: "Events service. Consumes events from NATS and persists them to ClickHouse " +
			"(raw, JS-executor, router, business-rule, trigger, workflow, processed/EVA, and DLQ " +
			"streams). This HTTP API serves the read side — cursor-paginated history queries over " +
			"each event stream, processed-event search with EVA dynamic-field filters, asset " +
			"connectivity history, and retention-policy management.",
		Contact: &swagger.Contact{
			Name: "Mapex Solutions",
			URL:  "https://mapexos.io",
		},
		License: &swagger.License{
			Name: "Business Source License 1.1",
			URL:  "https://github.com/Mapex-Solutions/mapexOS/blob/main/LICENSE",
		},
		ExternalDocs: &swagger.ExternalDocs{
			Description: "Project repository",
			URL:         "https://github.com/Mapex-Solutions/mapexOS",
		},
		Servers: []swagger.Server{
			{
				URL:         fmt.Sprintf("http://%s:%d", httpAddress, httpPort),
				Description: fmt.Sprintf("%s environment", goEnv),
			},
		},
		Tags: []swagger.Tag{
			{
				Name: "Events",
				Description: "Query the event history stored in ClickHouse — raw, JS-executor, router, " +
					"business-rule, trigger, workflow, processed (EVA), and DLQ streams. Scoped to the " +
					"caller's organization.",
			},
			{
				Name: "Asset Status",
				Description: "Query asset connectivity (online/offline) history derived from the event " +
					"streams.",
			},
			{
				Name: "Retention",
				Description: "Manage retention policies that govern how long each event stream is kept " +
					"in ClickHouse.",
			},
		},
	})
}

// BuildSwagger assembles the recorded routes into the OpenAPI document and serves
// it at /swagger. Call after all modules have registered their routes.
func BuildSwagger(app *web.App) {
	swagger.Build(app)
}
