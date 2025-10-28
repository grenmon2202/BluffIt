package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/grenmon2202/bluffit/internal/ws/handler"
	"github.com/grenmon2202/bluffit/internal/ws/requests"
	"github.com/grenmon2202/bluffit/internal/ws/response"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func sendResponse(ws *websocket.Conn, res any) bool {
	err := ws.WriteJSON(res)
	if err != nil {
		fmt.Println("write error:", err)
		return false
	}
	return true
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer ws.Close()

	for {
		var msg requests.Envelope
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("client closed:", err)
			} else {
				fmt.Println("read error:", err)
			}
			break
		}

		var res response.ResponseBody

		switch msg.Type {
		case "CreateLobby":
			var body requests.CreateLobbyRequest
			if err := json.Unmarshal(msg.Data, &body); err != nil {
				fmt.Println("unmarshal error:", err)
				continue
			}

			res = handler.CreateLobby(body)

		case "FetchLobby":
			var body requests.FetchLobbyRequest
			if err := json.Unmarshal(msg.Data, &body); err != nil {
				fmt.Println("unmarshal error:", err)
				continue
			}

			res = handler.FetchLobby(body)

		case "JoinLobby":
			var body requests.JoinLobbyRequest
			if err := json.Unmarshal(msg.Data, &body); err != nil {
				fmt.Println("unmarshal error:", err)
				continue
			}

			res = handler.JoinLobby(body)

		case "SendMessage":
			var body requests.SendMessageRequest
			if err := json.Unmarshal(msg.Data, &body); err != nil {
				fmt.Println("unmarshal error:", err)
				continue
			}

			res = handler.SendMessage(body)

		default:
			fmt.Println("Unknown message type:", msg.Type)

			res = response.ResponseBody{
				Message: "Unknown message type",
				Status:  400,
			}
		}

		if !sendResponse(ws, &res) {
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
