package cleanlinter

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
)

func Test_layer_String(t *testing.T) {
	tests := []struct {
		name     string
		l        layer
		expected string
	}{
		{
			name:     "domain",
			l:        layerDomain,
			expected: "Domain layer",
		},
		{
			name:     "usecase",
			l:        layerUseCase,
			expected: "UseCase layer",
		},
		{
			name:     "adapter",
			l:        layerAdapter,
			expected: "Adapter layer",
		},
		{
			name:     "infra",
			l:        layerInfra,
			expected: "Infrastructure layer",
		},
		{
			name:     "unknown",
			l:        100,
			expected: "Unknown layer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.l.String()

			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestLinter_CheckImports(t *testing.T) {
	type fields struct {
		pathToDomain  string
		pathToUseCase string
		pathToAdapter string
		pathToInfra   string
		verboseMode   bool
	}
	validPathUseCase1 := "github.com/bullgare/cleanlinter/internal/usecase/pkg/1"
	validPathUseCase2 := "github.com/bullgare/cleanlinter/internal/usecase/pkg/1"
	validPathDomain1 := "github.com/bullgare/cleanlinter/internal/domain/pkg/1"
	validPathDomain2 := "github.com/bullgare/cleanlinter/internal/domain/pkg/2"

	validFields := fields{
		pathToDomain:  "github.com/bullgare/cleanlinter/internal/domain",
		pathToUseCase: "github.com/bullgare/cleanlinter/internal/usecase",
		pathToAdapter: "github.com/bullgare/cleanlinter/internal/adapter",
		pathToInfra:   "github.com/bullgare/cleanlinter/internal/infra",
		verboseMode:   false,
	}

	tests := []struct {
		name        string
		fields      fields
		pass        *analysis.Pass
		expectedMsg string
		expectedErr error
	}{
		{
			name:   "happy path - usecase imports domain",
			fields: validFields,
			pass: &analysis.Pass{
				Pkg: types.NewPackage(validPathUseCase1, ""),
				Fset: func() *token.FileSet {
					fset := token.NewFileSet()
					fset.AddFile("usecase_filename.go", 111, 2)
					return fset
				}(),
				Files: []*ast.File{
					{
						Package: 111,
						Imports: []*ast.ImportSpec{
							{Path: &ast.BasicLit{Value: validPathDomain1}},
							{Path: &ast.BasicLit{Value: validPathDomain2}},
						},
					},
				},
			},
			expectedMsg: "",
			expectedErr: nil,
		},
		{
			name:   "error - domain imports usecase",
			fields: validFields,
			pass: &analysis.Pass{
				Pkg: types.NewPackage(validPathDomain1, ""),
				Fset: func() *token.FileSet {
					fset := token.NewFileSet()
					fset.AddFile("usecase_filename.go", 111, 2)
					return fset
				}(),
				Files: []*ast.File{
					{
						Package: 111,
						Imports: []*ast.ImportSpec{
							{Path: &ast.BasicLit{Value: validPathUseCase1, ValuePos: 1}},
							{Path: &ast.BasicLit{Value: validPathUseCase2, ValuePos: 2}},
						},
					},
				},
			},
			expectedMsg: "Domain layer is not allowed to reference UseCase layer",
			expectedErr: nil,
		},
		{
			name:   "not a clean architecture path - return nil, no checks",
			fields: validFields,
			pass: &analysis.Pass{
				Pkg: types.NewPackage("github.com/bullgare/cleanlinter/some/other/path", ""),
				Fset: func() *token.FileSet {
					fset := token.NewFileSet()
					fset.AddFile("usecase_filename.go", 111, 2)
					return fset
				}(),
				Files: []*ast.File{
					{
						Package: 111,
						Imports: []*ast.ImportSpec{
							{Path: &ast.BasicLit{Value: validPathUseCase1, ValuePos: 1}},
							{Path: &ast.BasicLit{Value: validPathUseCase2, ValuePos: 2}},
						},
					},
				},
			},
			expectedMsg: "",
			expectedErr: nil,
		},
		{
			name:        "no pass - return error",
			fields:      validFields,
			pass:        nil,
			expectedMsg: "",
			expectedErr: errors.New("unexpected nil pass"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Linter{
				PathToDomain:  tt.fields.pathToDomain,
				PathToUseCase: tt.fields.pathToUseCase,
				PathToAdapter: tt.fields.pathToAdapter,
				PathToInfra:   tt.fields.pathToInfra,
				verboseMode:   tt.fields.verboseMode,
			}

			if tt.pass != nil {
				tt.pass.Report = func(d analysis.Diagnostic) {
					assert.Equal(t, tt.expectedMsg, d.Message)
				}
			}

			err := l.CheckImports(tt.pass)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestLinter_getLayerByPackage(t *testing.T) {
	type fields struct {
		pathToDomain  string
		pathToUseCase string
		pathToAdapter string
		pathToInfra   string
		verboseMode   bool
	}

	validFields := fields{
		pathToDomain:  "github.com/bullgare/cleanlinter/internal/domain",
		pathToUseCase: "github.com/bullgare/cleanlinter/internal/usecase",
		pathToAdapter: "github.com/bullgare/cleanlinter/internal/adapter",
		pathToInfra:   "github.com/bullgare/cleanlinter/internal/infra",
	}

	tests := []struct {
		name     string
		fields   fields
		pkgPath  string
		expected layer
	}{
		{
			name:     "domain sub",
			fields:   validFields,
			pkgPath:  validFields.pathToDomain + "/a/pkg",
			expected: layerDomain,
		},
		{
			name:     "usecase sub",
			fields:   validFields,
			pkgPath:  validFields.pathToUseCase + "/a/pkg",
			expected: layerUseCase,
		},
		{
			name:     "adapter sub",
			fields:   validFields,
			pkgPath:  validFields.pathToAdapter + "/a/pkg",
			expected: layerAdapter,
		},
		{
			name:     "infra sub",
			fields:   validFields,
			pkgPath:  validFields.pathToInfra + "/a/pkg",
			expected: layerInfra,
		},
		{
			name:     "not from the project",
			fields:   validFields,
			pkgPath:  "/a/pkg",
			expected: layerNone,
		},
		{
			name:     "from the project, not listed",
			fields:   validFields,
			pkgPath:  "github.com/bullgare/cleanlinter/internal",
			expected: layerNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Linter{
				PathToDomain:  tt.fields.pathToDomain,
				PathToUseCase: tt.fields.pathToUseCase,
				PathToAdapter: tt.fields.pathToAdapter,
				PathToInfra:   tt.fields.pathToInfra,
				verboseMode:   tt.fields.verboseMode,
			}

			res := l.getLayerByPackage(tt.pkgPath)

			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestLinter_checkPrefix(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		prefix   string
		expected bool
	}{
		{
			name:     "happy path - full match",
			path:     "github.com/bullgare/cleanlinter/internal/domain",
			prefix:   "github.com/bullgare/cleanlinter/internal/domain",
			expected: true,
		},
		{
			name:     "happy path - prefix match",
			path:     "github.com/bullgare/cleanlinter/internal/domain/a/pkg",
			prefix:   "github.com/bullgare/cleanlinter/internal/domain",
			expected: true,
		},
		{
			name:     "subpath does not match exactly - no match",
			path:     "github.com/bullgare/cleanlinter/internal/domain123/a/pkg",
			prefix:   "github.com/bullgare/cleanlinter/internal/domain",
			expected: false,
		},
		{
			name:     "empty prefix does not match anything",
			path:     "github.com/bullgare/cleanlinter/internal/domain/a/pkg",
			prefix:   "",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Linter{}.checkPrefix(tt.path, tt.prefix)

			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestLinter_checkLayerReference(t *testing.T) {
	tests := []struct {
		name        string
		curLayer    layer
		importLayer layer
		expected    string
	}{
		{
			name:        "happy path, infra to unknown",
			curLayer:    layerInfra,
			importLayer: layerNone,
			expected:    "",
		},
		{
			name:        "happy path, infra to domain",
			curLayer:    layerInfra,
			importLayer: layerDomain,
			expected:    "",
		},
		{
			name:        "happy path, infra to usecase",
			curLayer:    layerInfra,
			importLayer: layerUseCase,
			expected:    "",
		},
		{
			name:        "happy path, infra to adapter",
			curLayer:    layerInfra,
			importLayer: layerAdapter,
			expected:    "",
		},
		{
			name:        "happy path, infra to infra",
			curLayer:    layerInfra,
			importLayer: layerInfra,
			expected:    "",
		},
		{
			name:        "happy path, adapter to unknown",
			curLayer:    layerAdapter,
			importLayer: layerNone,
			expected:    "",
		},
		{
			name:        "happy path, adapter to domain",
			curLayer:    layerAdapter,
			importLayer: layerDomain,
			expected:    "",
		},
		{
			name:        "happy path, adapter to usecase",
			curLayer:    layerAdapter,
			importLayer: layerUseCase,
			expected:    "",
		},
		{
			name:        "happy path, adapter to adapter",
			curLayer:    layerAdapter,
			importLayer: layerAdapter,
			expected:    "",
		},
		{
			name:        "adapter to infra is not allowed",
			curLayer:    layerAdapter,
			importLayer: layerInfra,
			expected:    "Adapter layer is not allowed to reference Infrastructure layer",
		},
		{
			name:        "happy path, usecase to unknown",
			curLayer:    layerUseCase,
			importLayer: layerNone,
			expected:    "",
		},
		{
			name:        "happy path, usecase to domain",
			curLayer:    layerUseCase,
			importLayer: layerDomain,
			expected:    "",
		},
		{
			name:        "happy path, usecase to usecase",
			curLayer:    layerUseCase,
			importLayer: layerUseCase,
			expected:    "",
		},
		{
			name:        "usecase to adapter is not allowed",
			curLayer:    layerUseCase,
			importLayer: layerAdapter,
			expected:    "UseCase layer is not allowed to reference Adapter layer",
		},
		{
			name:        "usecase to infra is not allowed",
			curLayer:    layerUseCase,
			importLayer: layerInfra,
			expected:    "UseCase layer is not allowed to reference Infrastructure layer",
		},
		{
			name:        "happy path, domain to unknown",
			curLayer:    layerDomain,
			importLayer: layerNone,
			expected:    "",
		},
		{
			name:        "happy path, domain to domain",
			curLayer:    layerDomain,
			importLayer: layerDomain,
			expected:    "",
		},
		{
			name:        "domain to usecase is not allowed",
			curLayer:    layerDomain,
			importLayer: layerUseCase,
			expected:    "Domain layer is not allowed to reference UseCase layer",
		},
		{
			name:        "domain to adapter is not allowed",
			curLayer:    layerDomain,
			importLayer: layerAdapter,
			expected:    "Domain layer is not allowed to reference Adapter layer",
		},
		{
			name:        "domain to infra is not allowed",
			curLayer:    layerDomain,
			importLayer: layerInfra,
			expected:    "Domain layer is not allowed to reference Infrastructure layer",
		},
		{
			name:        "happy path, unknown to unknown",
			curLayer:    layerNone,
			importLayer: layerNone,
			expected:    "",
		},
		{
			name:        "happy path, unknown to domain",
			curLayer:    layerNone,
			importLayer: layerDomain,
			expected:    "",
		},
		{
			name:        "happy path, unknown to usecase",
			curLayer:    layerNone,
			importLayer: layerUseCase,
			expected:    "",
		},
		{
			name:        "happy path, unknown to adapter",
			curLayer:    layerNone,
			importLayer: layerAdapter,
			expected:    "",
		},
		{
			name:        "happy path, unknown to infra",
			curLayer:    layerNone,
			importLayer: layerInfra,
			expected:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Linter{}.checkLayerReference(tt.curLayer, tt.importLayer)

			assert.Equal(t, tt.expected, err)
		})
	}
}
