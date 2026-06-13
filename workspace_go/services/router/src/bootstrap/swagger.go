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
		Description: "Router service. Consumes events from NATS, resolves asset and template " +
			"context via a tiered cache, evaluates RouteGroup match rules, and fans events out " +
			"to downstream subjects (events store, lakehouse, notifications, rule engine, " +
			"triggers, workflow). This HTTP API manages the RouteGroups that drive that routing.",
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
				Name: "Route Groups",
				Description: "Create, query, and manage RouteGroups — the match rules and routers " +
					"that decide where each incoming event is delivered. Scoped to the caller's organization.",
			},
		},
	})
}

// BuildSwagger assembles the recorded routes into the OpenAPI document and serves
// it at /swagger. Call after all modules have registered their routes.
func BuildSwagger(app *web.App) {
	swagger.Build(app)
}
