package cleanlinter

import (
	"flag"
	"fmt"

	"golang.org/x/tools/go/analysis"
)

var flagSet flag.FlagSet

var (
	pathToDomain  string
	pathToUseCase string
	pathToAdapter string
	pathToInfra   string

	verboseMode bool
)

func init() {
	flagSet.StringVar(&pathToDomain, "cleanlinter_path_to_domain", "", "path to Domain layer; usually contains domain layer business rules")
	flagSet.StringVar(&pathToUseCase, "cleanlinter_path_to_usecase", "", "path to UseCase layer; usually contains services, reusable handlers, the app core logic goes here")
	flagSet.StringVar(&pathToAdapter, "cleanlinter_path_to_adapter", "", "path to Adapter layer; usually contains adapters to other APIs")
	flagSet.StringVar(&pathToInfra, "cleanlinter_path_to_infra", "", "path to Infrastructure layer; usually contains presentation logic (http-, grpc-handlers), DB integration, and so on")
	flagSet.BoolVar(&verboseMode, "cleanlinter_verbose", false, "verbose mode")
}

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:  "cleanlinter",
		Doc:   "Clean Architecture linter. Checks your project's dependencies according to Clean Architecture pattern",
		Flags: flagSet,
		Run:   run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	countNonEmptyPaths := 0
	if pathToDomain != "" {
		countNonEmptyPaths++
	}
	if pathToUseCase != "" {
		countNonEmptyPaths++
	}
	if pathToAdapter != "" {
		countNonEmptyPaths++
	}
	if pathToInfra != "" {
		countNonEmptyPaths++
	}

	if countNonEmptyPaths < 2 {
		return nil, fmt.Errorf("please specify at least 2 layers. Otherwise there is nothing to check (got %d)", countNonEmptyPaths)
	}

	linter := Linter{
		PathToDomain:  pathToDomain,
		PathToUseCase: pathToUseCase,
		PathToAdapter: pathToAdapter,
		PathToInfra:   pathToInfra,
		verboseMode:   verboseMode,
	}
	err := linter.CheckImports(pass)

	return nil, err
}
