package services

import (
	"encoding/json"
	"github.com/digital-technology-agency/secret-santa/pkg/data/bitcask"
	"github.com/digital-technology-agency/secret-santa/pkg/models"
	"math/rand"
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
	players := []*models.Player{}
	friends := []*models.Player{}
	for _, bytes := range allData {
		var playerData models.Player
		err = json.Unmarshal(bytes, &playerData)
		if err != nil {
			return err
		}
		players = append(players, &playerData)
		friends = append(friends, &playerData)
	}
	markerData := map[string]*models.Player{}
	for len(players) != 0 && len(friends) != 0 {
		playersLen := len(players)
		friendsLen := len(friends)
		indexPlayer := rand.Intn(playersLen)
		indexFriend := rand.Intn(friendsLen)
		selectPlayer := players[indexPlayer]
		selectFriend := friends[indexFriend]
		if selectPlayer.Id == selectFriend.Id {
			continue
		}
		if markerData[selectPlayer.Id] != nil {
			players = append(players[:indexPlayer], players[indexPlayer+1:]...)
			continue
		}
		if markerData[selectFriend.Id] != nil {
			friends = append(friends[:indexFriend], friends[indexFriend+1:]...)
			continue
		}
		markerData[selectPlayer.Id] = &models.Player{
			Id:       selectPlayer.Id,
			Login:    selectPlayer.Login,
			FriendId: selectFriend.Id,
		}
		markerData[selectFriend.Id] = &models.Player{
			Id:       selectFriend.Id,
			Login:    selectFriend.Login,
			FriendId: selectPlayer.Id,
		}
		friends = append(friends[:indexFriend], friends[indexFriend+1:]...)
		players = append(players[:indexPlayer], players[indexPlayer+1:]...)
	}
	for _, updatePlayer := range markerData {
		err := game.AddPlayer(*updatePlayer)
		if err != nil {
			return err
		}
	}
	defer func() {
		markerData = nil
		players = nil
		friends = nil
	}()
	return nil
}
