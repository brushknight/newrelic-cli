package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var outputFormat string
var outputPlain bool

//const defaultProfileName string = "default"

// Command represents the base command when called without any subcommands
var Command = &cobra.Command{
	PersistentPreRun:  initializeCLI,
	Use:               appName,
	Short:             "The New Relic CLI",
	Long:              `The New Relic CLI enables users to perform tasks against the New Relic APIs`,
	Version:           version,
	DisableAutoGenTag: true, // Do not print generation date on documentation
}

func initializeCLI(cmd *cobra.Command, args []string) {
	if client.Client == nil {
		client.Client = createClient()
	}
}

func createClient() *newrelic.NewRelic {
	c, err := client.NewClient(configuration.GetActiveProfileName())
	if err != nil {
		// An error was encountered initializing the client.  This may not be a
		// problem since many commands don't require the use of an initialized client
		log.Debugf("error initializing client: %s", err)
	}

	return c
}

// func initializeProfile(ctx context.Context) {
// 	var accountID int
// 	var region string
// 	var licenseKey string
// 	var insightsInsertKey string
// 	var err error

// 	if c.DefaultProfile != "" {
// 		err = errors.New("default profile already exists, not attempting to initialize")
// 		return
// 	}

// 	apiKey := os.Getenv("NEW_RELIC_API_KEY")
// 	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
// 	region = os.Getenv("NEW_RELIC_REGION")
// 	licenseKey = os.Getenv("NEW_RELIC_LICENSE_KEY")
// 	insightsInsertKey = os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")

// 	// If we don't have a personal API key we can't initialize a profile.
// 	if apiKey == "" {
// 		err = errors.New("api key not provided, not attempting to initialize default profile")
// 		return
// 	}

// 	// Default the region to US if it's not in the environment
// 	if region == "" {
// 		region = "US"
// 	}

// 	// Use the accountID from the environment if we have it.
// 	if envAccountID != "" {
// 		accountID, err = strconv.Atoi(envAccountID)
// 		if err != nil {
// 			err = fmt.Errorf("couldn't parse account ID: %s", err)
// 			return
// 		}
// 	}

// 	// We should have an API key by this point, initialize the client.
// 	client := createClient()

// 	// If we still don't have an account ID try to look one up from the API.
// 	if accountID == 0 {
// 		accountID, err = fetchAccountID(client)
// 		if err != nil {
// 			return
// 		}
// 	}

// 	if licenseKey == "" {
// 		// We should have an account ID by now, so fetch the license key for it.
// 		licenseKey, err = fetchLicenseKey(ctx, client, accountID)
// 		if err != nil {
// 			log.Error(err)
// 			return
// 		}
// 	}

// 	if insightsInsertKey == "" {
// 		// We should have an API key by now, so fetch the insights insert key for it.
// 		insightsInsertKey, err = fetchInsightsInsertKey(client, accountID)
// 		if err != nil {
// 			log.Error(err)
// 		}
// 	}

// 	if !hasProfileWithDefaultName(c.Profiles) {
// 		p := credentials.Profile{
// 			Region:            region,
// 			APIKey:            apiKey,
// 			AccountID:         accountID,
// 			LicenseKey:        licenseKey,
// 			InsightsInsertKey: insightsInsertKey,
// 		}

// 		err = c.AddProfile(defaultProfileName, p)
// 		if err != nil {
// 			return
// 		}

// 		log.Infof("profile %s added", text.FgCyan.Sprint(defaultProfileName))
// 	}

// 	if len(c.Profiles) == 1 {
// 		err = c.SetDefaultProfile(defaultProfileName)
// 		if err != nil {
// 			err = fmt.Errorf("error setting %s as the default profile: %s", text.FgCyan.Sprint(defaultProfileName), err)
// 			return
// 		}

// 		log.Infof("setting %s as default profile", text.FgCyan.Sprint(defaultProfileName))
// 	}

// 	if err != nil {
// 		log.Debugf("couldn't initialize default profile: %s", err)
// 	}
// }

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	Command.Use = appName
	Command.Version = version
	Command.SilenceUsage = os.Getenv("CI") != ""

	// Silence Cobra's internal handling of error messaging
	// since we have a custom error handler in main.go
	Command.SilenceErrors = true

	return Command.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	Command.PersistentFlags().StringVar(&outputFormat, "format", output.DefaultFormat.String(), "output text format ["+output.FormatOptions()+"]")
	Command.PersistentFlags().BoolVar(&outputPlain, "plain", false, "output compact text")
}

func initConfig() {
	utils.LogIfError(output.SetFormat(output.ParseFormat(outputFormat)))
	utils.LogIfError(output.SetPrettyPrint(!outputPlain))
}
