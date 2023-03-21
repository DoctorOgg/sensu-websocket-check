// websockets.go
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func showHelp() {
	fmt.Println("Usage: " + os.Args[0] + " [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	listenPort := flag.Int("port", 8080, "port to listen on")
	listenAddress := flag.String("address", "0.0.0.0", "address to listen on")
	helpflag := flag.Bool("help", false, "show help")

	flag.Parse()

	if *helpflag {
		showHelp()
		os.Exit(0)
	}

	fmt.Println("Starting server on ", *listenAddress, ":", *listenPort)

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is a websocket echo server. take a look at /echo")
	})

	http.ListenAndServe(*listenAddress+":"+strconv.Itoa(*listenPort), nil)
}
