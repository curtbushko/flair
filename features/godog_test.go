// Package features provides BDD acceptance tests for flair using godog.
// These tests implement the Gherkin specifications from PLAN.md and validate
// the complete flair pipeline from domain layer through adapters to CLI.
package features

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"github.com/curtbushko/flair/features/steps"
)

// opts holds godog options configured for the test run.
var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
	Paths:  []string{"domain", "adapters", "application", "e2e"},
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: steps.InitializeScenario,
		Options:             &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("BDD tests failed")
	}
}
