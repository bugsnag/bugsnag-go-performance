package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var scenariosMap = map[string] func()(string, func()){
	"HandledScenario": HandledScenario,
}

func HandledScenario() (string, func()) {
	f := func () {
		fmt.Println("HELLO WORLD")
	}
	return "OUTPUT", f
}

func main() {
		fmt.Println("[Bugsnag] Starting testapp")
		// Listening to the OS Signals
		ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		ticker := time.NewTicker(1 * time.Second)

		addr := os.Getenv("DEFAULT_MAZE_ADDRESS")
		if (addr == "") {
			addr = DEFAULT_MAZE_ADDRESS
		}
	
		for {
			select {
			case <-ticker.C:
				fmt.Println("[Bugsnag] Get command")
				command := GetCommand(addr)
				fmt.Printf("[Bugsnag] Received command: %+v\n", command)

				if command.Action == "run-scenario" {
					prepareScenarioFunc, ok := scenariosMap[command.ScenarioName]
					if ok {
						_, scenarioFunc := prepareScenarioFunc()
						//bugsnag.Configure(config)
						scenarioFunc()
					}
				}
			case <-ctx.Done():
					fmt.Println("[Bugsnag] Context is done, closing")
					ticker.Stop()
					return
			}
		}
}
