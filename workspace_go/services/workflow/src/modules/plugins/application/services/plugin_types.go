package services

import (
	"workflow/src/modules/plugins/application/di"
)

type PluginService struct {
	deps di.PluginServiceDependenciesInjection
}
