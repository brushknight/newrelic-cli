package execution

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/pkg/region"
)

type SuccessLinkGenerator interface {
	GenerateExplorerLink(filter string) string
	GenerateEntityLink(entityGUID string) string
	GenerateRedirectURL(status InstallStatus) string
}

type ConcreteSuccessLinkGenerator struct{}

var nrPlatformHostnames = struct {
	Staging string
	US      string
	EU      string
}{
	Staging: "staging-one.newrelic.com",
	US:      "one.newrelic.com",
	EU:      "one.eu.newrelic.com",
}

func NewConcreteSuccessLinkGenerator() *ConcreteSuccessLinkGenerator {
	return &ConcreteSuccessLinkGenerator{}
}

func (g *ConcreteSuccessLinkGenerator) GenerateExplorerLink(filter string) string {
	return generateExplorerLink(filter)
}

func (g *ConcreteSuccessLinkGenerator) GenerateEntityLink(entityGUID string) string {
	return generateEntityLink(entityGUID)
}

func toJSON(data interface{}) string {
	c, _ := json.MarshalIndent(data, "", "  ")

	return string(c)
}

// GenerateRedirectURL creates a URL for the user to navigate to after running
// through an installation. The URL is displayed in the CLI out as well and is
// also provided in the nerdstorage document. This provides the user two options
// to see their data - click from the CLI output or from the frontend.
func (g *ConcreteSuccessLinkGenerator) GenerateRedirectURL(status InstallStatus) string {
	log.Print("\n\n **************************** \n")
	log.Printf("\n GenerateRedirectURL - status:       %+v \n", status)
	log.Printf("\n GenerateRedirectURL - status json:  %+v \n", toJSON(status))

	if status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED) {
		switch t := status.successLinkConfig.Type; {
		case strings.EqualFold(string(t), "explorer"):
			return g.GenerateExplorerLink(status.successLinkConfig.Filter)
		default:
			return g.GenerateEntityLink(status.HostEntityGUID())
		}
	}

	return ""
}

func generateExplorerLink(filter string) string {
	log.Printf("\n generateExplorerLink - filter:  %+v \n", filter)
	log.Print("\n **************************** \n\n")

	return fmt.Sprintf("https://%s/launcher/nr1-core.explorer?platform[filters]=%s&platform[accountId]=%d",
		nrPlatformHostname(),
		utils.Base64Encode(filter),
		credentials.DefaultProfile().AccountID,
	)
}

func generateEntityLink(entityGUID string) string {
	log.Printf("\n generateEntityLink - entityGUID:  %+v \n", entityGUID)
	log.Print("\n **************************** \n\n")

	return fmt.Sprintf("https://%s/redirect/entity/%s", nrPlatformHostname(), entityGUID)
}

// nrPlatformHostname returns the host for the platform based on the region set.
func nrPlatformHostname() string {
	defaultProfile := credentials.DefaultProfile()
	if defaultProfile == nil {
		return nrPlatformHostnames.US
	}

	if strings.EqualFold(defaultProfile.Region, region.Staging.String()) {
		return nrPlatformHostnames.Staging
	}

	if strings.EqualFold(defaultProfile.Region, region.EU.String()) {
		return nrPlatformHostnames.EU
	}

	return nrPlatformHostnames.US
}
