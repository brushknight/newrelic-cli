package execution

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestInstallEventsReporter_RecipeFailed(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeFailed(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)

}

func TestInstallEventsReporter_RecipeInstalling(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeInstalling(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeInstalled(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeSkipped(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeSkipped(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeRecommended(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeRecommended(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_writeStatus(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")

	recipes := []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			DisplayName:    types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	createInstallEventCallCount := 0

	err := r.RecipesAvailable(status, recipes)
	createInstallEventCallCount++
	require.NoError(t, err)
	require.Equal(t, createInstallEventCallCount, c.CreateInstallEventCallCount)

	err = r.RecipesSelected(status, recipes)
	createInstallEventCallCount++
	require.NoError(t, err)
	require.Equal(t, createInstallEventCallCount, c.CreateInstallEventCallCount)

	manifest := types.DiscoveryManifest{}

	err = r.DiscoveryComplete(status, manifest)
	createInstallEventCallCount++
	require.NoError(t, err)
	require.Equal(t, createInstallEventCallCount, c.CreateInstallEventCallCount)

	for _, testRecipe := range recipes {
		err = r.RecipeAvailable(status, testRecipe)
		createInstallEventCallCount++
		require.NoError(t, err)
		require.Equal(t, createInstallEventCallCount, c.CreateInstallEventCallCount)
	}

	err = r.InstallComplete(status)
	createInstallEventCallCount++
	require.NoError(t, err)
	require.Equal(t, createInstallEventCallCount, c.CreateInstallEventCallCount)

	err = r.InstallCanceled(status)
	createInstallEventCallCount++
	require.NoError(t, err)
	require.Equal(t, createInstallEventCallCount, c.CreateInstallEventCallCount)

}
