# Vanilla Runtime
This is the most trivial example of a runtime. It only provides a single source and sink, both of which are NATS.
The purpose of this runtime is to show how you can create your own runtimes, nothing more.

> [!IMPORTANT]
> THIS RUNTIME IS A SIMPLE EXAMPLE TO SHOW HOW TO CREATE A RUNTIME. IT IS NOT INTENDED FOR PRODUCTION USE.

## Crafting a Runtime
If you want to write your own runtime, you need to be aware of how things work. A runtime is the component which will
take a connector definition and run it. An agent is responsible for launching the runtime and passing the connector
configuration to it. Capturing logs, events and metrics is the responsibility of the agent, not the runtime itself. 
This allows the runtime to focus on the task of running the connector and not worry about the rest.

As any go application, you will need to have a `main.go` file which will be the entry point of your runtime. This file
will be responsible for setting up the runtime and starting it. The `github.com/synadia-io/connect` module provides 
some convenience functions to help you with this:

```go
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
```

The `runtime.FromEnv()` function will return a runtime which is configured from the environment. This is useful as it
allows you to configure your runtime using environment variables. The `rt.Launch()` function will start the runtime and
hand it the workload function. The workload function is responsible for the actual work of the runtime. In this case,
the `workload.Run` function is used.

The majority of work when writing a runtime will be in the workload function.

## Workload Function

```go 
func(ctx context.Context, runtime *runtime.Runtime, v model.ConnectorConfig) error
```

A workload function will perform the actual work based on the connector configuration provided. You can implement this
function in any way you see fit. The `model.ConnectorConfig` struct is a struct which represents the connector
configuration.

## .vent Directory
Every runtime must have a `.vent` directory. This directory contains meta information about the runtime as well as the
components provided by the runtime. 

### runtime.yml
The `.vent` directory must contain a `runtime.yml` file which describes the runtime:

```yaml
model_version: "1"
name: vanilla
label: Vanilla
version: 0.0.1
image: ghcr.io/synadia-io/connect-runtime-vanilla:latest
author:
  name: Synadia
  email: code@synadia.com
description: |-
  The vanilla runtime serves as an example of what a simple runtime might look like. It provides only a single input
  and output, and is intended to be used as a starting point for creating more complex runtimes.
```

The `model_version` dictates the format of the file and must be set to `1`. The `name` is the name of the runtime while
the `label` provides a human readable label. The version is a SemVer version of the runtime used for versioning.

Runtimes are linked to an image. This image is the image which will be used to execute the runtime. A prefix notes 
which kind of image it is: 

- `docker:ghcr.io/synadia-io/connect-runtime-vanilla:latest` - A docker image
- `nex:IMAGES/vanilla/latest` - A Nex image

If no prefix is provided, the image is assumed to be a docker image.

The author section provides information about the author of the runtime (you). This makes it easier for people to reach
out in case they have questions or issues. Or if they just want to say thank you.

The description is a human readable description of what this runtime is all about. Don't hold back, tell people what
makes your runtime special.

### sinks
The `sinks` directory contains the sink components provided by the runtime. Each sink has its own definition file which
describes the configuration options and overall structure.

### sources
The `sources` directory contains the source components provided by the runtime. Similar to a sink, each source has 
its own definition file.

## Contributing
We love to get feedback, bug reports, and contributions from our community. If you have any questions or want to
contribute, feel free to reach out in the #connectors channel on the NATS slack.

Take a look at the [Contributing](CONTRIBUTING.md) guide to get started.