package adapter

import (
	"github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/infra" // want "Adapter layer is not allowed to reference Infrastructure layer"
)

var AdapterAPI = infra.Implementation()
