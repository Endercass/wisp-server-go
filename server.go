package main

import (
	"fmt"
	"net/http"

	wispparse "github.com/Endercass/wisp-server-go/pkg/wisp-parse"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	var cp wispparse.ConnectPacket = wispparse.ConnectPacket{
		StreamType:          wispparse.StreamType(1),
		DestinationPort:     8080,
		DestinationHostname: "localhost",
	}

	fmt.Println(cp.Marshal())

	fmt.Println("Server starting on port 8080")
	http.ListenAndServe(":8080", nil)
}
