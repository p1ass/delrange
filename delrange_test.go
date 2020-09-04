package delrange_test

import (
	"testing"

	"github.com/p1ass/delrange"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, delrange.Analyzer, "a")
}
