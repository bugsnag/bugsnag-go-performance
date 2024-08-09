package main

import (
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
		signalsChan := make(chan os.Signal, 1)
		signal.Notify(signalsChan, syscall.SIGINT, syscall.SIGTERM)
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
			case <-signalsChan:
					fmt.Println("[Bugsnag] Signal received, closing")
					ticker.Stop()
					return
			}
		}
}
