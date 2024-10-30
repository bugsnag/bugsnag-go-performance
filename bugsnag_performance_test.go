package bugsnagperformance

import (
	"log"
	"os"
)

func resetEnv() {
	os.Clearenv()
	Config = Configuration{
		ReleaseStage: "production",
		Logger:       log.New(os.Stdout, "[BugsnagPerformance] ", log.LstdFlags),
	}
}
