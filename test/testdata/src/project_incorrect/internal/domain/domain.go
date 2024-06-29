package domain

import "github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/usecase" // want "Domain layer is not allowed to reference UseCase layer"

var Entity = usecase.Service
