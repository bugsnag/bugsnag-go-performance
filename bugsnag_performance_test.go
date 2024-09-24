package bugsnagperformance

import (
	"os"
)

func resetEnv() {
	os.Clearenv()
	Config = Configuration{
		ReleaseStage: "production",
	}
}
