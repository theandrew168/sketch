package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"

	"github.com/coder/websocket"
)

//go:embed index.html
var indexHTML []byte

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(indexHTML)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Accepting new websocket connection")

		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			fmt.Println("Error accepting websocket connection:", err)
			return
		}
		defer c.CloseNow()

		fmt.Println("Accepted websocket connection")

		ctx := context.Background()
		for {
			kind, msg, err := c.Read(ctx)
			if err != nil {
				fmt.Println("Error reading websocket message:", err)
				return
			}

			fmt.Printf("Received message of kind %v: %s\n", kind, string(msg))

			err = c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				fmt.Println("Error writing websocket message:", err)
				return
			}
		}
	})

	addr := "127.0.0.1:8080"
	fmt.Println("Listening on http://" + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
