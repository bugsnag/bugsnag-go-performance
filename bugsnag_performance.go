package bugsnagperformance

import (
	"fmt"
	"sync"
)

// Version defines the version of this Bugsnag performance module
const Version = "0.1.0"

// Config is the configuration for the default Bugsnag performance module
var Config Configuration

var readEnvConfigOnce sync.Once

func init() {
	fmt.Println("Starting bugsnag performance")
	Config = Configuration{
		ReleaseStage: "production",
	}
}

// Configure Bugsnag. The only required setting is the APIKey, which can be
// obtained by clicking on "Settings" in your Bugsnag dashboard.
// Returns OTeL sampler, trace exporter and error
func Configure(config Configuration) (interface{}, interface{}, error) {

	readEnvConfigOnce.Do(Config.loadEnv)
	Config.update(&config)
	err := Config.validate()
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}
