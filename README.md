[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/DoctorOgg/sensu-websocket-check)

# sensu-websocket-check

## Table of Contents

- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The sensu-websocket-check is a [Sensu Check][6] that is used to check a WebSocket server. you can specify a string to check for in the response, and a payload to send to the server.  The plugin will exit with a 0 if the string is found in the response, and a 2 if it is not found.

## Files

- bin/
  - sensu-websocket-check - The main executable for the plugin

## Usage examples

```bash
./sensu-websocket-check 
Usage:
  sensu-go-websocket-check [flags]
  sensu-go-websocket-check [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -h, --help                     help for sensu-go-websocket-check
  -i, --ignore-cert              Ignore certificate errors
  -p, --payload string           Payload to send to the WebSocket server (default "ping")
  -s, --string-to-check string   String to check in the response (default "ping")
  -t, --timeout int              Timeout in seconds (default 10)
  -u, --url string               URL of the WebSocket server to check (e.g., ws://example.com/socket))

Use "sensu-go-websocket-check [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add DoctorOgg/sensu-websocket-check
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/DoctorOgg/sensu-websocket-check].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-websocket-check
  namespace: default
spec:
  command: sensu-websocket-check -u ws://localhost:8080/echo -s ping -p pings
  subscriptions:
  - system
  runtime_assets:
  - DoctorOgg/sensu-websocket-check
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-websocket-check repository:

```
go build
```

## Testing the plugin locally

there is a test server in test-server directory.  To run it, do the following:

```bash
go run test-server/main.go
```

Note this server also supports some command line options:

```bash
$ go run test-server/main.go -help

Options:
  -address string
        address to listen on (default "0.0.0.0")
  -help
        show help
  -port int
        port to listen on (default 8080)
```

and then run the plugin with the following command:

```bash
$ ./sensu-websocket-check -u ws://localhost:8080/echo -s ping -p pings
sending payload: pings to ws://localhost:8080/echo
unexpected response: pings%

$ echo $?
2
```
