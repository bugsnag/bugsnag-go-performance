package bugsnagperformance

import (
	"context"
	"log"
	"os"
)

func resetEnv() {
	os.Clearenv()
	Config = Configuration{
		ReleaseStage: "production",
		Logger:       log.New(os.Stdout, "[BugsnagPerformance] ", log.LstdFlags),
		MainContext:  context.TODO(),
		Endpoint:     "",
	}
}
