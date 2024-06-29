package cleanlinter_test

import (
	"testing"

	"github.com/bullgare/cleanlinter"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestIntegration(t *testing.T) {
	tests := []struct {
		name            string
		packagesToCheck []string
		flags           map[string]string
	}{
		{
			name: "project_correct",
			packagesToCheck: []string{
				"github.com/bullgare/cleanlinter/test/testdata/src/project_correct",
				"github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/domain",
				"github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/usecase",
				"github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/adapter",
				"github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/infra",
			},
			flags: map[string]string{
				"cleanlinter_path_to_domain":  "github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/domain",
				"cleanlinter_path_to_usecase": "github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/usecase",
				"cleanlinter_path_to_adapter": "github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/adapter",
				"cleanlinter_path_to_infra":   "github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/infra",
				"cleanlinter_verbose":         "false",
			},
		},
		{
			name: "project_incorrect",
			packagesToCheck: []string{
				"github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect",
				"github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/domain",
				"github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/usecase",
				"github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/adapter",
				"github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/infra",
			},
			flags: map[string]string{
				"cleanlinter_path_to_domain":  "github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/domain",
				"cleanlinter_path_to_usecase": "github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/usecase",
				"cleanlinter_path_to_adapter": "github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/adapter",
				"cleanlinter_path_to_infra":   "github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect/internal/infra",
				"cleanlinter_verbose":         "false",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := cleanlinter.NewAnalyzer()
			for k, v := range tt.flags {
				analyzer.Flags.Set(k, v)
			}

			analysistest.RunWithSuggestedFixes(
				t,
				analysistest.TestData(),
				analyzer,
				tt.packagesToCheck...,
			)
		})
	}
}
