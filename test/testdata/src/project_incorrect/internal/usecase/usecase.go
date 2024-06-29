package usecase

import "github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/adapter" // want "UseCase layer is not allowed to reference Adapter layer"

var Service = adapter.AdapterAPI
