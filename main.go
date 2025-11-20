// TODO: Move request handling logic to broadcast/client.
// Add subscription request, only allow subsequent requests to flow through if subscribed (Except CreateLobby and JoinLobby)

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/grenmon2202/bluffit/internal/lobby"
	"github.com/grenmon2202/bluffit/ws/broadcast"
	"github.com/grenmon2202/bluffit/ws/handler"
	"github.com/grenmon2202/bluffit/ws/schema/request"
	"github.com/grenmon2202/bluffit/ws/schema/response"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var broadcast_hub = broadcast.CreateBroadcastHub()

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := broadcast.CreateBroadcastClient(ws)
	go client.WritePump()

	var lobbyID uint32
	defer func() {
		if lobbyID != 0 {
			broadcast_hub.RemoveSubscriber(lobbyID, client)
		}
		ws.Close()
	}()

	for {
		var msg request.Envelope
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("client closed:", err)
			} else {
				fmt.Println("read error:", err)
			}
			break
		}

		if msg.RequestID == "" {
			fmt.Println("Missing RequestID in message")
			res := response.ResponseBody{
				Message: "Missing RequestID in message",
				Status:  400,
			}
			client.Send(&res)
			continue
		}

		var res response.ResponseBody

		switch msg.Type {
		case "CreateLobby":
			var body request.CreateLobbyRequest
			if err := json.Unmarshal(msg.Data, &body); err != nil {
				fmt.Println("unmarshal error:", err)
				continue
			}

			res = handler.CreateLobby(body)

		case "FetchLobby": // this is also used as a handshake
			var body request.FetchLobbyRequest
			if err := json.Unmarshal(msg.Data, &body); err != nil {
				fmt.Println("unmarshal error:", err)
				continue
			}
			var fetchedLobby lobby.Lobby
			res, fetchedLobby = handler.FetchLobby(body)

			lobbyID = fetchedLobby.ID

			if res.Status == 200 {
				broadcast_hub.AddSubscriber(lobbyID, client)
			}

		case "JoinLobby":
			var body request.JoinLobbyRequest
			if err := json.Unmarshal(msg.Data, &body); err != nil {
				fmt.Println("unmarshal error:", err)
				continue
			}
			var joinedLobby lobby.Lobby
			res, joinedLobby = handler.JoinLobby(body)

			lobbyID = joinedLobby.ID

			if res.Status == 200 {
				broadcast_hub.AddSubscriber(lobbyID, client)
				broadcast_hub.BroadcastToLobby(lobbyID, response.BroadcastBody{
					Type:    "LobbyInfo",
					Message: "New player joined the lobby",
					Data:    res.Data.(map[string]any)["lobby"],
				})
			}

		case "LeaveLobby":
			var body request.LeaveLobbyRequest
			if err := json.Unmarshal(msg.Data, &body); err != nil {
				fmt.Println("unmarshal error:", err)
				continue
			}
			var leftLobby lobby.Lobby
			res, leftLobby = handler.LeaveLobby(body)

			if res.Status == 200 && lobbyID != 0 {
				broadcast_hub.BroadcastToLobby(lobbyID, response.BroadcastBody{
					Type:    "LobbyInfo",
					Message: "A player left the lobby",
					Data:    leftLobby,
				})
				broadcast_hub.RemoveSubscriber(lobbyID, client)
				return
			}

		default:
			fmt.Println("Unknown message type:", msg.Type)

			res = response.ResponseBody{
				Message: "Unknown message type",
				Status:  400,
			}
		}

		res.RequestID = msg.RequestID

		client.Send(&res)
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
