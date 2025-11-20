package handler

import (
	"fmt"

	"github.com/grenmon2202/bluffit/internal/lobby"
	"github.com/grenmon2202/bluffit/internal/player"
	"github.com/grenmon2202/bluffit/ws/schema/request"
	"github.com/grenmon2202/bluffit/ws/schema/response"
)

var globalLobbyStore = lobby.NewLobbyStore()

func CreateLobby(data request.CreateLobbyRequest) response.ResponseBody {
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

func FetchLobby(data request.FetchLobbyRequest) (response.ResponseBody, lobby.Lobby) {
	lobby, exists := globalLobbyStore.GetLobbyByID(data.LobbyID)
	if !exists {
		return response.ResponseBody{
			Message: "Lobby not found",
			Status:  404,
		}, *lobby
	}

	return response.ResponseBody{
		Message: "Lobby fetched successfully",
		Status:  200,
		Data:    *lobby,
	}, *lobby
}

func JoinLobby(data request.JoinLobbyRequest) (response.ResponseBody, lobby.Lobby) {
	success, message, player, lobby := globalLobbyStore.JoinLobbyByCode(data.LobbyCode, data.PlayerName)

	if !success {
		return response.ResponseBody{
			Message: message,
			Status:  400,
		}, lobby
	}

	response_data := map[string]any{
		"player": player,
		"lobby":  lobby,
	}

	return response.ResponseBody{
		Message: "Successfully joined lobby",
		Status:  200,
		Data:    response_data,
	}, lobby
}

func LeaveLobby(data request.LeaveLobbyRequest) (response.ResponseBody, lobby.Lobby) {
	lobby, exists := globalLobbyStore.GetLobbyByID(data.LobbyID)
	if !exists {
		return response.ResponseBody{
			Message: "Lobby not found",
			Status:  404,
		}, *lobby
	}

	if !lobby.Players.CheckPlayerExistsByID(data.PlayerID) {
		return response.ResponseBody{
			Message: "Player not found in the lobby",
			Status:  404,
		}, *lobby
	}

	lobby.Players.RemovePlayerByID(data.PlayerID)

	if lobby.LobbyHostID == data.PlayerID {
		globalLobbyStore.RemoveLobbyByID(lobby.ID)
		return response.ResponseBody{
			Message: "Host has left the lobby. Lobby closed.",
			Status:  200,
		}, *lobby
	}

	return response.ResponseBody{
		Message: "Successfully left the lobby",
		Status:  200,
	}, *lobby
}
