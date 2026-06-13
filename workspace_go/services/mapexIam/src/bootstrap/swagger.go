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
		Description: "Identity & Access Management service. Owns the core platform identities and " +
			"their relationships: users, organizations (with hierarchical org tree), memberships, " +
			"roles and permissions, groups, and reusable lists. Issues and refreshes auth tokens " +
			"and resolves per-request authorization coverage. The onboarding orchestrator composes " +
			"these to provision a new tenant in one call.",
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
			{Name: "Auth", Description: "Authenticate users and refresh access tokens."},
			{Name: "Users", Description: "Create, query, and manage platform users."},
			{Name: "Organizations", Description: "Manage organizations and the hierarchical org tree. Scoped to the caller's organization."},
			{Name: "Memberships", Description: "Manage the membership relation between users and organizations (with roles)."},
			{Name: "Me", Description: "Read and manage the authenticated caller's own profile, memberships, and context."},
			{Name: "Roles", Description: "Manage roles and their permission sets."},
			{Name: "Groups", Description: "Manage groups used to scope assets and apply policies."},
			{Name: "Lists", Description: "Manage reusable lists (key/value option sets) referenced across the platform."},
			{Name: "Onboarding", Description: "Provision a new tenant — organization, admin user, and baseline roles — in a single orchestrated call."},
		},
	})
}

// BuildSwagger assembles the recorded routes into the OpenAPI document and serves
// it at /swagger. Call after all modules have registered their routes.
func BuildSwagger(app *web.App) {
	swagger.Build(app)
}
