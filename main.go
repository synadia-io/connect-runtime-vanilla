package main

import (
	"context"
	"fmt"
	"github.com/synadia-io/connect-runtime-vanilla/workload"
	"github.com/synadia-io/connect/runtime"
	"os"
)

func main() {
	// A runtime is invoked with a single argument; a base64 encoded connector configuration.
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("usage: vanilla <config>")
		os.Exit(1)
	}

	// A runtime is constructed from the environment. This will return an error if the required environment variables
	// are not set.
	rt, err := runtime.FromEnv()
	if err != nil {
		panic(err)
	}

	// When launching a runtime, you will need to provide a workload function. This function will be invoked with the
	// runtime, and the connector configuration. The workload function is responsible for the actual work of your runtime.
	if err := rt.Launch(context.Background(), workload.Run, args[0]); err != nil {
		panic(err)
	}
}
