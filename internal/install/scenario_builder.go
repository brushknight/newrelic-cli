package install

import (
	"github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

type TestScenario string

const (
	Basic         TestScenario = "BASIC"
	Fail          TestScenario = "FAIL"
	ExecDiscovery TestScenario = "EXEC_DISCOVERY"
)

var (
	TestScenarios = []TestScenario{
		Basic,
		Fail,
	}
	emptyResults = []nrdb.NRDBResult{
		map[string]interface{}{
			"count": 0.0,
		},
	}
	nonEmptyResults = []nrdb.NRDBResult{
		map[string]interface{}{
			"count": 1.0,
		},
	}
)

func TestScenarioValues() []string {
	v := make([]string, len(TestScenarios))
	for i, s := range TestScenarios {
		v[i] = string(s)
	}

	return v
}

type ScenarioBuilder struct {
	installerContext types.InstallerContext
}

func NewScenarioBuilder(ic types.InstallerContext) *ScenarioBuilder {
	b := ScenarioBuilder{
		installerContext: ic,
	}

	return &b
}

func (b *ScenarioBuilder) BuildScenario(s TestScenario) *RecipeInstaller {
	switch s {
	case Basic:
		return b.Basic()
	case Fail:
		return b.Fail()
	case ExecDiscovery:
		return b.ExecDiscovery()
	}

	return nil
}

func (b *ScenarioBuilder) Basic() *RecipeInstaller {

	// mock implementations
	rf := setupRecipeFetcherGuidedInstall()
	ers := []execution.StatusSubscriber{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
	slg := execution.NewPlatformLinkGenerator()
	statusRollup := execution.NewInstallStatus(ers, slg)
	c := validation.NewMockNRDBClient()
	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)
	v := validation.NewPollingRecipeValidator(c)
	cv := diagnose.NewMockConfigValidator()
	mv := discovery.NewEmptyManifestValidator()
	lkf := NewMockLicenseKeyFetcher()
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer()
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	s := ux.NewPlainProgress()
	rfi := recipes.NewRecipeFilterRunner(b.installerContext, statusRollup)
	rvp := execution.NewRecipeVarProvider()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: s,
		configValidator:   cv,
		manifestValidator: mv,
		licenseKeyFetcher: lkf,
		recipeFilterer:    rfi,
		recipeVarPreparer: rvp,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func (b *ScenarioBuilder) Fail() *RecipeInstaller {

	// mock implementations
	rf := setupRecipeFetcherGuidedInstall()
	ers := []execution.StatusSubscriber{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
	slg := execution.NewPlatformLinkGenerator()
	statusRollup := execution.NewInstallStatus(ers, slg)
	c := validation.NewMockNRDBClient()
	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)
	v := validation.NewPollingRecipeValidator(c)
	cv := diagnose.NewMockConfigValidator()
	mv := discovery.NewEmptyManifestValidator()

	lkf := NewMockLicenseKeyFetcher()
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer()
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewMockFailingRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: pi,
		configValidator:   cv,
		manifestValidator: mv,
		licenseKeyFetcher: lkf,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func (b *ScenarioBuilder) ExecDiscovery() *RecipeInstaller {

	// mock implementations
	rf := setupRecipeFetcherExecDiscovery()
	ers := []execution.StatusSubscriber{
		execution.NewMockStatusReporter(),
		execution.NewTerminalStatusReporter(),
	}
	slg := execution.NewPlatformLinkGenerator()
	statusRollup := execution.NewInstallStatus(ers, slg)
	c := validation.NewMockNRDBClient()
	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)
	v := validation.NewPollingRecipeValidator(c)
	cv := diagnose.NewMockConfigValidator()
	mv := discovery.NewEmptyManifestValidator()

	lkf := NewMockLicenseKeyFetcher()
	ff := recipes.NewRecipeFileFetcher()
	d := discovery.NewPSUtilDiscoverer()
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewMockFailingRecipeExecutor()
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()
	rvp := execution.NewRecipeVarProvider()

	rr := recipes.NewRecipeFilterRunner(b.installerContext, statusRollup)

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: pi,
		configValidator:   cv,
		manifestValidator: mv,
		licenseKeyFetcher: lkf,
		recipeVarPreparer: rvp,
		recipeFilterer:    rr,
	}

	i.InstallerContext = b.installerContext

	return &i
}

func setupRecipeFetcherGuidedInstall() recipes.RecipeFetcher {
	f := recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:        "infrastructure-agent-installer",
			DisplayName: "Infrastructure Agent",
			PreInstall: types.OpenInstallationPreInstallConfiguration{
				Info: `
This is the Infrastructure Agent Installer preinstall message.
It is made up of a multi line string.
				`,
			},
			PostInstall: types.OpenInstallationPostInstallConfiguration{
				Info: `
This is the Infrastructure Agent Installer postinstall message.
It is made up of a multi line string.
				`,
			},
			ValidationNRQL: "test NRQL",
			Install: `
version: '3'
tasks:
  default:
`,
		},
		{
			Name:           "logs-integration",
			DisplayName:    "Logs integration",
			ValidationNRQL: "test NRQL",
			LogMatch: []types.OpenInstallationLogMatch{
				{
					Name: "docker log",
					File: "/var/lib/docker/containers/*/*.log",
				},
			},
			Install: `
version: '3'
tasks:
  default:
`,
		},
	}

	return f
}

func setupRecipeFetcherExecDiscovery() recipes.RecipeFetcher {
	f := recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:        "matching-recipe",
			DisplayName: "matching-recipe",
			PreInstall: types.OpenInstallationPreInstallConfiguration{
				RequireAtDiscovery: "true",
			},
		},
		{
			Name:        "non-matching-recipe",
			DisplayName: "non-matching-recipe",
			PreInstall: types.OpenInstallationPreInstallConfiguration{
				RequireAtDiscovery: "bogus command",
			},
		},
	}

	return f
}
