//+build integration

package configuration

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testCfg = `{
		"*":{
			"loglevel":"debug",
			"plugindir": "/Users/ctrombley/.newrelic/plugins",
			"prereleasefeatures": "NOT_ASKED",
			"sendusagedata": "NOT_ASKED",
			"testInt": 42,
			"testString": "value1",
			"teststring": "value2"
			"caseInsensitiveTest": "value"
		}
	}`
)

func TestConfigProvider_Ctor_NilOption(t *testing.T) {
	_, err := NewConfigProvider(nil)
	require.NoError(t, err)
}

func TestConfigProvider_Ctor_OptionError(t *testing.T) {
	_, err := NewConfigProvider(func(cp *ConfigProvider) error { return errors.New("") })
	require.Error(t, err)
}

func TestConfigProvider_Ctor_CaseInsensitiveKeyCollision(t *testing.T) {
	_, err := NewConfigProvider(
		WithFieldDefinitions(
			FieldDefinition{Key: "asdf"},
			FieldDefinition{Key: "ASDF"},
		),
	)
	require.Error(t, err)
}

func TestConfigProvider_Ctor_CaseSensitiveKeyOverlap(t *testing.T) {
	_, err := NewConfigProvider(
		WithFieldDefinitions(
			FieldDefinition{
				Key:           "asdf",
				CaseSensitive: true,
			},
		),
	)
	require.NoError(t, err)
}

func TestConfigProvider_GetString(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetString("loglevel")
	require.NoError(t, err)
	require.Equal(t, "debug", actual)
}

func TestConfigProvider_GetString_CaseSensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:           "testString",
			CaseSensitive: true,
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetString("testString")
	require.NoError(t, err)
	require.Equal(t, "value1", actual)

	actual, err = p.GetString("teststring")
	require.NoError(t, err)
	require.Equal(t, "value2", actual)
}

func TestConfigProvider_GetString_CaseInsensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:           "caseInsensitiveTest",
			CaseSensitive: false,
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetString("caseinsensitivetest")
	require.NoError(t, err)
	require.Equal(t, "value", actual)
}

func TestConfigProvider_GetString_NotDefined(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	_, err = p.GetString("undefined")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no value found")
}

func TestConfigProvider_GetString_DefaultValue(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:     "prereleasefeatures",
			Default: "NOT_ASKED",
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetString("prereleasefeatures")
	require.NoError(t, err)
	require.Equal(t, "NOT_ASKED", actual)
}

func TestConfigProvider_GetString_EnvVarOverride(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:     "prereleasefeatures",
			Default: "NOT_ASKED",
			EnvVar:  "NEW_RELIC_CLI_PRERELEASE",
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = os.Setenv("NEW_RELIC_CLI_PRERELEASE", "testValue")
	require.NoError(t, err)

	actual, err := p.GetString("prereleasefeatures")
	require.NoError(t, err)
	require.Equal(t, "testValue", actual)
}

func TestConfigProvider_GetInt(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetInt("testInt")
	require.NoError(t, err)
	require.Equal(t, int64(42), actual)
}

func TestConfigProvider_GetInt_NotDefined(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	_, err = p.GetInt("undefined")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no value found")
}

func TestConfigProvider_GetInt_DefaultValue(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:     "testInt",
			Default: 42,
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetInt("testInt")
	require.NoError(t, err)
	require.Equal(t, int64(42), actual)
}

func TestConfigProvider_GetInt_EnvVarOverride(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:    "testInt",
			EnvVar: "TEST_INT",
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = os.Setenv("TEST_INT", "42")
	require.NoError(t, err)

	actual, err := p.GetInt("testInt")
	require.NoError(t, err)
	require.Equal(t, int64(42), actual)
}

func TestConfigProvider_GetInt_EnvVarOverride_WrongType(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:    "testInt",
			EnvVar: "TEST_INT",
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = os.Setenv("TEST_INT", "TEST_VALUE")
	require.NoError(t, err)

	_, err = p.GetInt("testInt")
	require.Error(t, err)
}

func TestConfigProvider_Set(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.NoError(t, err)

	actual, err := p.GetString("loglevel")
	require.NoError(t, err)
	require.Equal(t, "trace", actual)
}

func TestConfigProvider_SetTernary(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("testTernary", TernaryValues.Allow)
	require.NoError(t, err)

	actual, err := p.GetTernary("testTernary")
	require.NoError(t, err)
	require.Equal(t, TernaryValues.Allow, actual)
}

func TestConfigProvider_SetTernary_Invalid(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:               "testTernary",
			SetValidationFunc: IsTernary(),
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("testTernary", Ternary("invalid"))
	require.Error(t, err)

	err = p.Set("anotherTestTernary", "invalid")
	require.NoError(t, err)

	actual, err := p.GetTernary("anotherTestTernary")
	require.NoError(t, err)
	require.False(t, actual.Bool())
}

func TestConfigProvider_Set_CaseSensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithExplicitValues(),
		WithFieldDefinitions(FieldDefinition{
			Key:           "loglevel",
			CaseSensitive: true,
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.NoError(t, err)

	err = p.Set("logLevel", "info")
	require.Error(t, err)
}

func TestConfigProvider_Set_CaseInsensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewConfigProvider(
		WithExplicitValues(),
		WithFieldDefinitions(FieldDefinition{
			Key:           "loglevel",
			CaseSensitive: false,
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.NoError(t, err)

	err = p.Set("LOGLEVEL", "info")
	require.NoError(t, err)

	actual, err := p.GetString("loglevel")
	require.NoError(t, err)
	require.Equal(t, "info", actual)
}

func TestConfigProvider_Set_FileDoesNotExist(t *testing.T) {
	p, err := NewConfigProvider(
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.NoError(t, err)

	actual, err := p.GetString("loglevel")
	require.NoError(t, err)
	require.Equal(t, "trace", actual)
}

func TestConfigProvider_Set_ExplicitValues_CaseInsensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithExplicitValues(),
		WithFieldDefinitions(FieldDefinition{
			Key: "allowed",
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.Error(t, err)

	err = p.Set("ALLOWED", "testValue")
	require.NoError(t, err)
}

func TestConfigProvider_Set_ExplicitValues_CaseSensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithExplicitValues(),
		WithFieldDefinitions(FieldDefinition{
			Key:           "allowed",
			CaseSensitive: true,
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.Error(t, err)

	err = p.Set("ALLOWED", "testValue")
	require.Error(t, err)

	err = p.Set("allowed", "testValue")
	require.NoError(t, err)
}

func TestConfigProvider_Set_ValidationFunc_IntGreaterThan(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:               "loglevel",
			SetValidationFunc: IntGreaterThan(0),
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", 0)
	require.Error(t, err)

	err = p.Set("loglevel", 1)
	require.NoError(t, err)
}

func TestConfigProvider_Set_ValidationFunc_IntGreaterThan_WrongType(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:               "loglevel",
			SetValidationFunc: IntGreaterThan(0),
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "debug")
	require.Error(t, err)

	err = p.Set("loglevel", 1)
	require.NoError(t, err)
}

func TestConfigProvider_Set_ValidationFunc_StringInStrings_CaseSensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:               "loglevel",
			SetValidationFunc: StringInStrings(true, "valid", "alsoValid"),
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.Error(t, err)

	err = p.Set("loglevel", "valid")
	require.NoError(t, err)
}

func TestConfigProvider_Set_ValidationFunc_StringInStrings_CaseInsensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:               "loglevel",
			SetValidationFunc: StringInStrings(false, "valid", "alsoValid"),
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "VALID")
	require.NoError(t, err)

	err = p.Set("loglevel", "ALSOVALID")
	require.NoError(t, err)
}

func TestConfigProvider_Set_ValidationFunc_StringInStrings_WrongType(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewConfigProvider(
		WithFieldDefinitions(FieldDefinition{
			Key:               "testInt",
			SetValidationFunc: StringInStrings(false, "valid", "alsoValid"),
		}),
		WithFilePersistence(f.Name()),
		WithScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("testInt", 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not a string")
}
