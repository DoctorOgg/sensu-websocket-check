package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

type Config struct {
	sensu.PluginConfig
	url           string
	StringToCheck string
	IgnoreCert    bool
	Payload       string
	Debug         bool
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
		&sensu.PluginConfigOption{
			Path:      "debug",
			Env:       "DEBUG",
			Argument:  "debug",
			Shorthand: "d",
			Default:   false,
			Usage:     "Enable debug mode",
			Value:     &plugin.Debug,
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
	headers := http.Header{}
	urlObj, _ := url.Parse(plugin.url)
	headers.Add("Host", urlObj.Host)
	headers.Add("User-Agent", "Sensu Go WebSocket Check")
	headers.Add("Origin", urlObj.Scheme+"://"+urlObj.Host)

	c, resp, err := websocket.DefaultDialer.Dial(plugin.url, headers)

	if plugin.Debug {
		fmt.Println("Response Status: " + resp.Status + "\n")
		fmt.Println("Response Headers:\n")
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
		fmt.Println("--------------------")

		// Print out the response body
		fmt.Println("Response Body:")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s\n", body)
		fmt.Println("--------------------")

	}

	if err != nil {
		fmt.Println("error during websocket connection: " + err.Error())
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
