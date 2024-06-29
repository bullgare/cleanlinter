package cleanlinter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type layer int

const (
	layerDomain layer = iota + 1
	layerUseCase
	layerAdapter
	layerInfra

	layerNone layer = -1
)

func (l layer) String() string {
	switch l {
	case layerDomain:
		return "Domain layer"
	case layerUseCase:
		return "UseCase layer"
	case layerAdapter:
		return "Adapter layer"
	case layerInfra:
		return "Infrastructure layer"
	default:
		return "Unknown layer"
	}
}

type Linter struct {
	PathToDomain  string
	PathToUseCase string
	PathToAdapter string
	PathToInfra   string

	verboseMode bool
}

func (l Linter) CheckImports(pass *analysis.Pass) error {
	b := &bytes.Buffer{}
	buf := bufio.NewWriter(b)
	defer func() {
		if l.verboseMode {
			buf.Flush()
			fmt.Println(b.String())
		}
	}()

	if pass == nil || pass.Pkg == nil {
		return errors.New("unexpected nil pass")
	}

	fmt.Fprintf(buf, "Running with params: pathToDomain=%s, pathToUseCase=%s, pathToAdapter=%s, pathToInfra=%s\n", l.PathToDomain, l.PathToUseCase, l.PathToAdapter, l.PathToInfra)

	curLayer := l.getLayerByPackage(pass.Pkg.Path())
	fmt.Fprintf(buf, "Checking package %s: %s\n", pass.Pkg.Path(), curLayer)
	if curLayer == layerNone {
		fmt.Fprintf(buf, "No need to check, returning\n")
		return nil
	}
	for _, f := range pass.Files {
		filePath := pass.Fset.Position(f.Package).Filename // @link https://stackoverflow.com/questions/72933175/go-get-filepath-from-ast-file

		fmt.Fprintf(buf, "File %s imports:\n", filePath)
		for _, imp := range f.Imports {
			impPath := imp.Path.Value
			impLayer := l.getLayerByPackage(impPath)
			msg := l.checkLayerReference(curLayer, impLayer)
			if msg != "" {
				pass.Reportf(imp.Pos(), msg)
			}
			fmt.Fprintf(buf, "	%s (%s)\n", impPath, impLayer)
		}
	}

	return nil
}

func (l Linter) getLayerByPackage(pkgPath string) layer {
	pkgPath = strings.Trim(pkgPath, "\"")
	if l.checkPrefix(pkgPath, l.PathToDomain) {
		return layerDomain
	}
	if l.checkPrefix(pkgPath, l.PathToUseCase) {
		return layerUseCase
	}
	if l.checkPrefix(pkgPath, l.PathToAdapter) {
		return layerAdapter
	}
	if l.checkPrefix(pkgPath, l.PathToInfra) {
		return layerInfra
	}
	return layerNone
}

func (l Linter) checkPrefix(path, prefix string) bool {
	if prefix == "" {
		return false
	}
	if path == prefix {
		return true
	}
	if strings.HasPrefix(path, prefix+"/") {
		return true
	}
	return false
}

func (l Linter) checkLayerReference(curLayer, importLayer layer) string {
	if curLayer == layerNone {
		return ""
	}

	if importLayer > curLayer {
		return fmt.Sprintf("%s is not allowed to reference %s", curLayer, importLayer)
	}
	return ""
}
