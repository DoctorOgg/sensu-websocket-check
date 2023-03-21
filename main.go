package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

type Config struct {
	sensu.PluginConfig
	url           string
	Timeout       int
	StringToCheck string
	IgnoreCert    bool
	Payload       string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-websocket-check",
			Short:    "Sensu check for WebSocket health",
			Keyspace: "sensu.io/plugins/sensu-go-websocket-check/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "url",
			Env:       "WEBSOCKET_URL",
			Argument:  "url",
			Shorthand: "u",
			Default:   "",
			Usage:     "URL of the WebSocket server to check (e.g., ws://example.com/socket))",
			Value:     &plugin.url,
		},
		&sensu.PluginConfigOption{
			Path:      "Timeout",
			Env:       "TIMEOUT",
			Argument:  "timeout",
			Shorthand: "t",
			Default:   10,
			Usage:     "Timeout in seconds",
			Value:     &plugin.Timeout,
		},
		&sensu.PluginConfigOption{
			Path:      "StringToCheck",
			Env:       "STRING_TO_CHECK",
			Argument:  "string-to-check",
			Shorthand: "s",
			Default:   "ping",
			Usage:     "String to check in the response",
			Value:     &plugin.StringToCheck,
		},
		&sensu.PluginConfigOption{
			Path:      "ignore-cert",
			Env:       "IGNORE_CERT",
			Argument:  "ignore-cert",
			Shorthand: "i",
			Default:   false,
			Usage:     "Ignore certificate errors",
			Value:     &plugin.IgnoreCert,
		},
		&sensu.PluginConfigOption{
			Path:      "payload",
			Env:       "PAYLOAD",
			Argument:  "payload",
			Shorthand: "p",
			Default:   "ping",
			Usage:     "Payload to send to the WebSocket server",
			Value:     &plugin.Payload,
		},
	}
)

func main() {
	useStdin := false
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if len(plugin.url) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--url or WEBSOCKET_URL environment variable is required")
	}
	if len(plugin.StringToCheck) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--string-to-check or STRING_TO_CHECK environment variable is required")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {

	if plugin.IgnoreCert {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	c, resp, err := websocket.DefaultDialer.Dial(plugin.url, http.Header{})

	if err == websocket.ErrBadHandshake {
		fmt.Printf("handshake failed with status %d", resp.StatusCode)
	}

	if err != nil {
		fmt.Println("error during websocket connection: " + err.Error())
		fmt.Println("response: " + resp.Status)
		fmt.Printf("error during websocket connection: %s", err.Error())
		return sensu.CheckStateCritical, nil
	}

	println("sending payload: " + plugin.Payload + " to " + plugin.url)
	err = c.WriteMessage(websocket.TextMessage, []byte(plugin.Payload))
	if err != nil {
		fmt.Printf("failed to write to WebSocket: %s", err)
		return sensu.CheckStateCritical, nil
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		fmt.Printf("failed to read from WebSocket: %s", err)
		return sensu.CheckStateCritical, nil
	}

	if string(message) != plugin.StringToCheck {
		fmt.Printf("unexpected response: %s", message)
		return sensu.CheckStateCritical, nil
	}
	c.Close()
	println("received message: " + string(message))

	return sensu.CheckStateOK, nil
}
