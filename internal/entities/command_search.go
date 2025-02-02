package entities

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

var cmdEntitySearch = &cobra.Command{
	Use:   "search",
	Short: "Search for New Relic entities",
	Long: `Search for New Relic entities

The search command performs a search for New Relic entities.
`,
	Example: "newrelic entity search --name <applicationName>",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			params := entities.EntitySearchQueryBuilder{}

			if entityName == "" && entityType == "" && entityAlertSeverity == "" && entityDomain == "" {
				utils.LogIfError(cmd.Help())
				log.Fatal("one of --name, --type, --alert-severity, or --domain are required")
			}

			if entityName != "" {
				params.Name = entityName
			}

			if entityType != "" {
				params.Type = entities.EntitySearchQueryBuilderType(entityType)
			}

			if entityAlertSeverity != "" {
				params.AlertSeverity = entities.EntityAlertSeverity(entityAlertSeverity)
			}

			if entityDomain != "" {
				params.Domain = entities.EntitySearchQueryBuilderDomain(entityDomain)
			}

			if entityTag != "" {
				key, value, err := assembleTagValue(entityTag)
				utils.LogIfFatal(err)

				params.Tags = []entities.EntitySearchQueryBuilderTag{{Key: key, Value: value}}
			}

			if entityReporting != "" {
				reporting, err := strconv.ParseBool(entityReporting)

				if err != nil {
					log.Fatalf("invalid value provided for flag --reporting. Must be true or false.")
				}

				params.Reporting = reporting
			}

			results, err := nrClient.Entities.GetEntitySearchWithContext(
				utils.SignalCtx,
				entities.EntitySearchOptions{},
				"",
				params,
				[]entities.EntitySearchSortCriteria{},
			)
			utils.LogIfFatal(err)

			entities := results.Results.Entities

			var result interface{}

			if len(entityFields) > 0 {
				mapped := mapEntities(entities, entityFields, utils.StructToMap)

				if len(mapped) == 1 {
					result = mapped[0]
				} else {
					result = mapped
				}
			} else {
				if len(entities) == 1 {
					result = entities[0]
				} else {
					result = entities
				}
			}

			utils.LogIfFatal(output.Print(result))
		})
	},
}

func mapEntities(entities []entities.EntityOutlineInterface, fields []string, fn utils.StructToMapCallback) []map[string]interface{} {
	mappedEntities := make([]map[string]interface{}, len(entities))

	for i, v := range entities {
		mappedEntities[i] = fn(v, fields)
	}

	return mappedEntities
}

func init() {
	Command.AddCommand(cmdEntitySearch)
	cmdEntitySearch.Flags().StringVarP(&entityName, "name", "n", "", "search for entities matching the given name")
	cmdEntitySearch.Flags().StringVarP(&entityType, "type", "t", "", "search for entities matching the given type")
	cmdEntitySearch.Flags().StringVarP(&entityAlertSeverity, "alert-severity", "a", "", "search for entities matching the given alert severity type")
	cmdEntitySearch.Flags().StringVarP(&entityReporting, "reporting", "r", "", "search for entities based on whether or not an entity is reporting (true or false)")
	cmdEntitySearch.Flags().StringVarP(&entityDomain, "domain", "d", "", "search for entities matching the given entity domain")
	cmdEntitySearch.Flags().StringVar(&entityTag, "tag", "", "search for entities matching the given entity tag")
	cmdEntitySearch.Flags().StringSliceVarP(&entityFields, "fields-filter", "f", []string{}, "filter search results to only return certain fields for each search result")
}
