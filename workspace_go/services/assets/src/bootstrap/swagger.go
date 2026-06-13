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
		Description: "Assets service. Owns the lifecycle of IoT Assets and their classification " +
			"templates, issues and manages MQTT client certificates, and tracks asset health. " +
			"It publishes a denormalized AssetReadModel to object storage that downstream " +
			"consumers (Router, JS-Executor, Events, the mapex-mqtt-broker plugin) hit as the " +
			"L3 fallback of their tiered caches. This HTTP API manages Assets, AssetTemplates, " +
			"and MQTT certificate issuance.",
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
				Name: "Assets",
				Description: "Create, query, and manage IoT Assets — the physical/logical devices " +
					"that emit events. Scoped to the caller's organization.",
			},
			{
				Name: "Asset Templates",
				Description: "Manage AssetTemplates — the classification and script definitions " +
					"(decode/validate/transform) shared by assets of the same kind.",
			},
			{
				Name: "MQTT Certificates",
				Description: "Issue and manage X.509 client certificates used by assets for " +
					"mutual-TLS MQTT authentication.",
			},
			{
				Name: "Gateway Certificates",
				Description: "Issue X.509 certificates for LoRaWAN gateway assets, reusing the " +
					"shared signer/CA machinery.",
			},
		},
	})
}

// BuildSwagger assembles the recorded routes into the OpenAPI document and serves
// it at /swagger. Call after all modules have registered their routes.
func BuildSwagger(app *web.App) {
	swagger.Build(app)
}
