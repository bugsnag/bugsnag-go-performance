package bugsnagperformance

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

// Needs to be first to test sync.Once for loadEnv
func TestConfigureTwiceEnv(t *testing.T) {
	resetEnv()
	testConfig := Configuration{
		APIKey:               "aaa",
		Endpoint:             "https://aaa.otlp.bugsnag.com/v1/traces",
		AppVersion:           "",
		ReleaseStage:         "prod1",
		EnabledReleaseStages: []string{"prod1", "prod2"},
		Logger:               nil}

	os.Setenv("BUGSNAG_PERFORMANCE_API_KEY", "aaa")
	os.Setenv("BUGSNAG_PERFORMANCE_RELEASE_STAGE", "prod1")
	os.Setenv("BUGSNAG_PERFORMANCE_ENABLED_RELEASE_STAGES", "prod1,prod2")

	_, err := Configure(Configuration{})
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testconfig is %+v\n", Config, testConfig)
	}

	os.Setenv("BUGSNAG_PERFORMANCE_API_KEY", "bbb")
	os.Setenv("BUGSNAG_PERFORMANCE_RELEASE_STAGE", "prod3")
	os.Setenv("BUGSNAG_PERFORMANCE_ENABLED_RELEASE_STAGES", "prod3,prod4")

	_, err = Configure(Configuration{APIKey: "aaa"})
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v\n", Config)
	}
}

func TestConfigureEmpty(t *testing.T) {
	resetEnv()
	_, err := Configure(Configuration{})
	if err == nil {
		t.Error("should return error on empty api key")
	}
}

func TestDefaultValues(t *testing.T) {
	resetEnv()
	testConfig := Configuration{
		APIKey:               "aaa",
		Endpoint:             "https://aaa.otlp.bugsnag.com/v1/traces",
		AppVersion:           "",
		ReleaseStage:         "production",
		EnabledReleaseStages: []string{},
		Logger:               nil}

	_, err := Configure(Configuration{APIKey: "aaa"})
	if err != nil {
		t.Error("should not return error")
	}
	// TODO prepare default logger
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testConfig is: %+v\n", Config, testConfig)
	}
}

func TestDefaultHubValues(t *testing.T) {
	resetEnv()
	testConfig := Configuration{
		APIKey:               "00000ffffeeee11112222333344445555",
		Endpoint:             "https://00000ffffeeee11112222333344445555.otlp.insighthub.smartbear.com/v1/traces",
		AppVersion:           "",
		ReleaseStage:         "production",
		EnabledReleaseStages: []string{},
		Logger:               nil}

	_, err := Configure(Configuration{
		APIKey: "00000ffffeeee11112222333344445555",
	})
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testConfig is: %+v\n", Config, testConfig)
	}
}

func TestConfigureOverwriteDefault(t *testing.T) {
	resetEnv()
	testConfig := Configuration{
		APIKey:               "bbb",
		Endpoint:             "myendpoint",
		AppVersion:           "123",
		ReleaseStage:         "dev",
		EnabledReleaseStages: []string{"dev"},
		Logger:               nil}

	_, err := Configure(testConfig)
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Error()
	}
}

func TestConfigureHubOverwriteDefault(t *testing.T) {
	resetEnv()
	testConfig := Configuration{
		APIKey:               "00000ffffeeee11112222333344445555",
		Endpoint:             "myendpoint",
		AppVersion:           "123",
		ReleaseStage:         "dev",
		EnabledReleaseStages: []string{"dev"},
		Logger:               nil}

	_, err := Configure(testConfig)
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Error()
	}
}

func TestConfigureMixedSetup(t *testing.T) {
	resetEnv()
	testConfig := Configuration{
		APIKey:               "bbb",
		Endpoint:             "https://bbb.otlp.bugsnag.com/v1/traces",
		AppVersion:           "123",
		ReleaseStage:         "prod1",
		EnabledReleaseStages: []string{},
		Logger:               nil}

	os.Setenv("BUGSNAG_API_KEY", "bbb")
	os.Setenv("BUGSNAG_APP_VERSION", "234")
	os.Setenv("BUGSNAG_PERFORMANCE_RELEASE_STAGE", "prod1")
	os.Setenv("BUGSNAG_RELEASE_STAGE", "prod2")

	// Has to be called manually, sync.Once already ran
	Config.loadEnv()
	_, err := Configure(Configuration{AppVersion: "123"})
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testconfig is %+v\n", Config, testConfig)
	}
}

func TestConfigureTwice(t *testing.T) {
	resetEnv()
	testConfig := Configuration{
		APIKey:               "aaa",
		Endpoint:             "https://aaa.otlp.bugsnag.com/v1/traces",
		AppVersion:           "",
		ReleaseStage:         "production",
		EnabledReleaseStages: []string{},
		Logger:               nil}

	_, err := Configure(testConfig)
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testconfig is %+v\n", Config, testConfig)
	}

	testConfig2 := Configuration{
		APIKey:               "bbb",
		Endpoint:             "https://bbb.otlp.bugsnag.com/v1/traces",
		AppVersion:           "",
		ReleaseStage:         "production2",
		EnabledReleaseStages: []string{},
		Logger:               nil}

	_, err = Configure(testConfig2)
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig2) {
		t.Errorf("Config is: %+v, testconfig2 is %+v\n", Config, testConfig2)
	}
}

func TestConfigureNotifierEnv(t *testing.T) {
	resetEnv()
	testConfig := Configuration{
		APIKey:               "aaa",
		Endpoint:             "https://aaa.otlp.bugsnag.com/v1/traces",
		AppVersion:           "version1",
		ReleaseStage:         "prod1",
		EnabledReleaseStages: []string{"prod1", "prod2"},
		Logger:               nil}

	os.Setenv("BUGSNAG_API_KEY", "aaa")
	os.Setenv("BUGSNAG_APP_VERSION", "version1")
	os.Setenv("BUGSNAG_RELEASE_STAGE", "prod1")
	os.Setenv("BUGSNAG_NOTIFY_RELEASE_STAGES", "prod1,prod2")

	// Has to be called manually, sync.Once already ran
	Config.loadEnv()
	_, err := Configure(Configuration{})
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testconfig is %+v\n", Config, testConfig)
	}
}

func TestEndpointFromEnvironment(t *testing.T) {
	customEndpoint := "https://endpoint.custom.com"
	hubAPIKey := "00000abcdef0123456789abcdef012345"
	setUp := func() {
		os.Setenv("BUGSNAG_PERFORMANCE_ENDPOINT", customEndpoint)
		os.Setenv("BUGSNAG_API_KEY", hubAPIKey)
	}

	t.Run("Should not override endpoint set by environment variable", func(st *testing.T) {
		resetEnv()
		setUp()

		testConfig := Configuration{
			Endpoint:     customEndpoint,
			APIKey:       hubAPIKey,
			ReleaseStage: "production",
		}

		// Has to be called manually, sync.Once already ran
		Config.loadEnv()
		_, err := Configure(Configuration{})
		if err != nil {
			t.Error("should not return error")
		}
		if !configsEqual(&Config, &testConfig) {
			t.Errorf("Config is: %+v, testconfig is %+v\n", Config, testConfig)
		}
	})

	t.Run("Should override endpoints set by environment variable with custom endpoint in code", func(st *testing.T) {
		resetEnv()
		setUp()
		c := &Configuration{}
		c.loadEnv()

		newCustomEndpoint := "https://test.endpoint.com"
		testConfig := Configuration{
			Endpoint:     newCustomEndpoint,
			APIKey:       hubAPIKey,
			ReleaseStage: "production",
		}

		// Has to be called manually, sync.Once already ran
		Config.loadEnv()
		_, err := Configure(Configuration{
			Endpoint: newCustomEndpoint,
		})
		if err != nil {
			t.Error("should not return error")
		}
		if !configsEqual(&Config, &testConfig) {
			t.Errorf("Config is: %+v, testconfig is %+v\n", Config, testConfig)
		}
	})

}

func configsEqual(first, second *Configuration) bool {
	enabledStagesFirst := len(first.EnabledReleaseStages)
	enabledStagesSecond := len(second.EnabledReleaseStages)

	if enabledStagesFirst == enabledStagesSecond {
		if enabledStagesFirst != 0 {
			sort.Strings(first.EnabledReleaseStages)
			sort.Strings(second.EnabledReleaseStages)
			if !reflect.DeepEqual(first.EnabledReleaseStages, second.EnabledReleaseStages) {
				return false
			}
		}
	} else {
		return false
	}

	return first.APIKey == second.APIKey &&
		first.AppVersion == second.AppVersion &&
		first.ReleaseStage == second.ReleaseStage &&
		first.Endpoint == second.Endpoint
}
