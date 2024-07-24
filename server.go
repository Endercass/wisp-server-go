package main

import (
	"fmt"
	"net/http"

	wispparse "github.com/Endercass/wisp-server-go/pkg/wisp-parse"
	"github.com/gorilla/websocket"
)

func handleWs(ws *websocket.Conn) {
	defer ws.Close()

	var initialContinue = wispparse.BuildContinuePacket(1024).ToPacket(0).Marshal()

	err := ws.WriteMessage(websocket.BinaryMessage, initialContinue)
	if err != nil {
		fmt.Println("Error writing message:", err)
		return
	}

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		packet, err := wispparse.ParsePacket(data)
		if err != nil {
			fmt.Println("Error parsing packet:", err)
			return
		}

		switch packet.Type {
		case wispparse.PacketTypeConnect:
			connectPacket, err := packet.ConnectPacket()
			if err != nil {
				fmt.Println("Error parsing connect packet:", err)
				return
			}
			fmt.Println("Connect packet:", connectPacket)

		case wispparse.PacketTypeData:
			dataPacket, err := packet.DataPacket()
			if err != nil {
				fmt.Println("Error parsing data packet:", err)
				return
			}
			fmt.Println("Data packet:", dataPacket)

		case wispparse.PacketTypeContinue:
			continuePacket, err := packet.ContinuePacket()
			if err != nil {
				fmt.Println("Error parsing continue packet:", err)
				return
			}
			fmt.Println("Continue packet:", continuePacket)

		case wispparse.PacketTypeClose:
			closePacket, err := packet.ClosePacket()
			if err != nil {
				fmt.Println("Error parsing close packet:", err)
				return
			}
			fmt.Println("Close packet:", closePacket)

		default:
			fmt.Println("Unknown packet type:", packet.Type)
		}

	}
}

func upgradeRoute(w http.ResponseWriter, r *http.Request) {
	// Determine if the request is a websocket upgrade request
	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Expected websocket connection", http.StatusBadRequest)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}

	handleWs(ws)
}

func main() {
	http.HandleFunc("/", upgradeRoute)

	var cp wispparse.ConnectPacket = wispparse.ConnectPacket{
		StreamType:          wispparse.StreamType(1),
		DestinationPort:     8080,
		DestinationHostname: "localhost",
	}

	fmt.Println("Server starting on port 8080")
	http.ListenAndServe(":8080", nil)
}
