package bugsnagperformance

import (
	"fmt"
	"os"
	"strings"
)

type Configuration struct {
	// Your Bugsnag API key, e.g. "c9d60ae4c7e70c4b6c4ebd3e8056d2b8". You can
	// find this by clicking Settings on https://bugsnag.com/.
	APIKey string

	// The currently running version of the app. This is used to filter errors
	// in the Bugsnag dasboard. If you set this then Bugsnag will only re-open
	// resolved errors if they happen in different app versions.
	AppVersion string

	// Address to which bugsnag will send traces to
	Endpoint string

	// The current release stage. This defaults to "production" and is used to
	// filter errors in the Bugsnag dashboard.
	// Should have the same value as "deployment.environment" otel resource attribute
	ReleaseStage string

	// The Release stages to send traces in. If you set this then bugsnag-go-performance will
	// only send traces to Bugsnag if the ReleaseStage is listed here.
	EnabledReleaseStages []string

	// TODO write logger - wrapper for OTEL logs
	Logger interface{}
}

func (config *Configuration) update(other *Configuration) *Configuration {
	if other.APIKey != "" {
		config.APIKey = other.APIKey
	}
	if other.AppVersion != "" {
		config.AppVersion = other.AppVersion
	}
	if other.Endpoint != "" {
		config.Endpoint = other.Endpoint
	}
	if other.ReleaseStage != "" {
		config.ReleaseStage = other.ReleaseStage
	}
	if other.EnabledReleaseStages != nil {
		config.EnabledReleaseStages = other.EnabledReleaseStages
	}
	if other.Logger != nil {
		config.Logger = other.Logger
	}

	return config
}

func (config *Configuration) validate() error {
	if config.APIKey == "" {
		return fmt.Errorf("no Bugsnag API Key set")
	}

	if config.Endpoint == "" {
		defaultEndpoint := fmt.Sprintf("https://%+v.otlp.bugsnag.com/v1/traces", config.APIKey)
		config.Endpoint = defaultEndpoint
	}

	return nil
}

func (config *Configuration) loadEnv() {
	envConfig := Configuration{}

	if apiKey := os.Getenv("BUGSNAG_PERFORMANCE_API_KEY"); apiKey != "" {
		envConfig.APIKey = apiKey
	} else if apiKey := os.Getenv("BUGSNAG_API_KEY"); apiKey != "" {
		envConfig.APIKey = apiKey
	}

	if appVersion := os.Getenv("BUGSNAG_APP_VERSION"); appVersion != "" {
		envConfig.AppVersion = appVersion
	}

	if stage := os.Getenv("BUGSNAG_PERFORMANCE_RELEASE_STAGE"); stage != "" {
		envConfig.ReleaseStage = stage
	} else if stage := os.Getenv("BUGSNAG_RELEASE_STAGE"); stage != "" {
		envConfig.ReleaseStage = stage
	}

	if stages := os.Getenv("BUGSNAG_PERFORMANCE_ENABLED_RELEASE_STAGES"); stages != "" {
		envConfig.EnabledReleaseStages = strings.Split(stages, ",")
	} else if stages := os.Getenv("BUGSNAG_NOTIFY_RELEASE_STAGES"); stages != "" {
		envConfig.EnabledReleaseStages = strings.Split(stages, ",")
	}

	config.update(&envConfig)
}