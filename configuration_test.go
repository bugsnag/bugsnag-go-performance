package bugsnagperformance

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func resetEnv() {
	os.Clearenv()
	Config = Configuration{
		ReleaseStage: "production",
	}
}

// Needs to be first to test sync.Once for loadEnv
func TestConfigureTwiceEnv(t *testing.T) {
	t.Cleanup(resetEnv)

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

	_, _, err := Configure(Configuration{})
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testconfig is %+v\n", Config, testConfig)
	}

	os.Setenv("BUGSNAG_PERFORMANCE_API_KEY", "bbb")
	os.Setenv("BUGSNAG_PERFORMANCE_RELEASE_STAGE", "prod3")
	os.Setenv("BUGSNAG_PERFORMANCE_ENABLED_RELEASE_STAGES", "prod3,prod4")

	_, _, err = Configure(Configuration{APIKey: "aaa"})
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v\n", Config)
	}
}

func TestConfigureEmpty(t *testing.T) {
	t.Cleanup(resetEnv)

	_, _, err := Configure(Configuration{})
	if err == nil {
		t.Error("should return error on empty api key")
	}
}

func TestDefaultValues(t *testing.T) {
	t.Cleanup(resetEnv)

	testConfig := Configuration{
		APIKey:               "aaa",
		Endpoint:             "https://aaa.otlp.bugsnag.com/v1/traces",
		AppVersion:           "",
		ReleaseStage:         "production",
		EnabledReleaseStages: []string{},
		Logger:               nil}

	_, _, err := Configure(Configuration{APIKey: "aaa"})
	if err != nil {
		t.Error("should not return error")
	}
	// TODO prepare default logger
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testConfig is: %+v\n", Config, testConfig)
	}
}

func TestConfigureOverwriteDefault(t *testing.T) {
	t.Cleanup(resetEnv)

	testConfig := Configuration{
		APIKey:               "bbb",
		Endpoint:             "myendpoint",
		AppVersion:           "123",
		ReleaseStage:         "dev",
		EnabledReleaseStages: []string{"dev"},
		Logger:               nil}

	_, _, err := Configure(testConfig)
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Error()
	}
}

func TestConfigureMixedSetup(t *testing.T) {
	t.Cleanup(resetEnv)

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
	_, _, err := Configure(Configuration{AppVersion: "123"})
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testconfig is %+v\n", Config, testConfig)
	}
}

func TestConfigureTwice(t *testing.T) {
	t.Cleanup(resetEnv)

	testConfig := Configuration{
		APIKey:               "aaa",
		Endpoint:             "https://aaa.otlp.bugsnag.com/v1/traces",
		AppVersion:           "",
		ReleaseStage:         "production",
		EnabledReleaseStages: []string{},
		Logger:               nil}

	_, _, err := Configure(testConfig)
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

	_, _, err = Configure(testConfig2)
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig2) {
		t.Errorf("Config is: %+v, testconfig2 is %+v\n", Config, testConfig2)
	}
}

func TestConfigureNotifierEnv(t *testing.T) {
	t.Cleanup(resetEnv)

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
	_, _, err := Configure(Configuration{})
	if err != nil {
		t.Error("should not return error")
	}
	if !configsEqual(&Config, &testConfig) {
		t.Errorf("Config is: %+v, testconfig is %+v\n", Config, testConfig)
	}
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
