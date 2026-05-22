package services

import (
	"workflow/src/modules/archiver/application/di"
)

type ArchiverService struct {
	deps di.ArchiverServiceDependenciesInjection
}
