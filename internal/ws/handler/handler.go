package handler

import (
	"fmt"

	"github.com/grenmon2202/bluffit/internal/lobby"
	"github.com/grenmon2202/bluffit/internal/player"
	"github.com/grenmon2202/bluffit/internal/ws/requests"
	"github.com/grenmon2202/bluffit/internal/ws/response"
)

var globalLobbyStore = lobby.NewLobbyStore()

func CreateLobby(data requests.CreateLobbyRequest) response.ResponseBody {
	hostPlayer := player.CreatePlayer(data.LobbyHost, true)

	lobby := globalLobbyStore.AddLobby(hostPlayer)
	fmt.Println("Created lobby:", lobby.ID)

	response_data := map[string]any{
		"lobby":  lobby,
		"player": hostPlayer,
	}

	return response.ResponseBody{
		Message: "Lobby created successfully",
		Status:  201,
		Data:    response_data,
	}
}

func FetchLobby(data requests.FetchLobbyRequest) response.ResponseBody {
	lobby, exists := globalLobbyStore.GetLobbyByID(data.LobbyID)
	if !exists {
		return response.ResponseBody{
			Message: "Lobby not found",
			Status:  404,
		}
	}

	return response.ResponseBody{
		Message: "Lobby fetched successfully",
		Status:  200,
		Data:    *lobby,
	}
}

func JoinLobby(data requests.JoinLobbyRequest) response.ResponseBody {
	success, message, player := globalLobbyStore.JoinLobbyByCode(data.LobbyCode, data.PlayerName)

	if !success {
		return response.ResponseBody{
			Message: message,
			Status:  400,
		}
	}

	response_data := map[string]any{
		"player": player,
	}

	return response.ResponseBody{
		Message: "Successfully joined lobby",
		Status:  200,
		Data:    response_data,
	}
}

func SendMessage(data requests.SendMessageRequest) response.ResponseBody {
	lobby, exists := globalLobbyStore.GetLobbyByID(data.LobbyID)
	if !exists {
		return response.ResponseBody{
			Message: "Lobby not found",
			Status:  404,
		}
	}

	success, message := lobby.AddChatMessage(data.Content, data.SenderID)
	if !success {
		return response.ResponseBody{
			Message: message,
			Status:  400,
		}
	}

	return response.ResponseBody{
		Message: "Message sent successfully",
		Status:  200,
	}
}
