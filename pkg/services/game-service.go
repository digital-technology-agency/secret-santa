package services

import (
	"encoding/json"
	"github.com/digital-technology-agency/secret-santa/pkg/data/bitcask"
	"github.com/digital-technology-agency/secret-santa/pkg/models"
)

type Game struct {
	data    *bitcask.Data
	ChatId  string
	Players []models.Player
}

//GetOrCreate get chat data
func GetOrCreate(chatId string) (*Game, error) {
	connect, err := bitcask.Connect(chatId)
	if err != nil {
		return nil, err
	}
	/*init data*/
	allData, err := connect.GetAll()
	if err != nil {
		return nil, err
	}
	gameResult := &Game{
		data:    connect,
		ChatId:  chatId,
		Players: []models.Player{},
	}
	for _, bytes := range allData {
		var playerData models.Player
		err = json.Unmarshal(bytes, &playerData)
		if err != nil {
			return nil, err
		}
		gameResult.Players = append(gameResult.Players, playerData)
	}
	return gameResult, nil
}

//AddPlayer add player to game
func (game *Game) AddPlayer(player models.Player) error {
	key := []byte(player.Id)
	value, err := json.Marshal(player)
	if err != nil {
		return err
	}
	return game.data.Add(key, value)
}

//Algorithm random
func (game *Game) Algorithm() error {
	/*init data*/
	allData, err := game.data.GetAll()
	if err != nil {
		return err
	}
	players := []models.Player{}
	for _, bytes := range allData {
		var playerData models.Player
		err = json.Unmarshal(bytes, &playerData)
		if err != nil {
			return err
		}
		players = append(players, playerData)
	}
	markerData := map[string]*models.Player{}
	for _, player := range players {
		if markerData[player.Id] != nil {
			continue
		}
		markerData[player.Id] = &models.Player{
			Id:       player.Id,
			Login:    player.Login,
			FriendId: player.FriendId,
		}
	}
	return nil
}
