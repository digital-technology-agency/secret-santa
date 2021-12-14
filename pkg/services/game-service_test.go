package services

import (
	"fmt"
	"github.com/digital-technology-agency/secret-santa/pkg/models"
	"testing"
)

var testChatId = "test-chat-id"

func TestGetOrCreate(t *testing.T) {
	type args struct {
		chatId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "get or create chat data",
			args: args{
				chatId: testChatId,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetOrCreate(tt.args.chatId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGame_AddPlayer(t *testing.T) {
	gameService, err := GetOrCreate(testChatId)
	if err != nil {
		t.Errorf("GetOrCreate() error = %v", err)
	}
	type args struct {
		player models.Player
	}
	tests := []struct {
		name    string
		fields  *Game
		args    args
		wantErr bool
	}{
		{
			name:   "add player to data",
			fields: gameService,
			args: args{
				player: models.Player{
					Id:       "1",
					Login:    "@testUserTg",
					FriendId: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Game{
				data:    tt.fields.data,
				ChatId:  tt.fields.ChatId,
				Players: tt.fields.Players,
			}
			if err := game.AddPlayer(tt.args.player); (err != nil) != tt.wantErr {
				t.Errorf("AddPlayer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_Algorithm(t *testing.T) {
	gameService, err := GetOrCreate(testChatId)
	if err != nil {
		t.Errorf("GetOrCreate() error = %v", err)
	}
	for i := 1; i < 11; i++ {
		err = gameService.AddPlayer(models.Player{
			Id:       fmt.Sprintf("%d", i),
			Login:    fmt.Sprintf("@user-%d", i),
			FriendId: "",
		})
		if err != nil {
			t.Errorf("AddPlayer() error = %v", err)
		}
	}
	tests := []struct {
		name    string
		fields  *Game
		wantErr bool
	}{
		{
			name:    "algorithm",
			fields:  gameService,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := &Game{
				data:    tt.fields.data,
				ChatId:  tt.fields.ChatId,
				Players: tt.fields.Players,
			}
			if err := game.Algorithm(); (err != nil) != tt.wantErr {
				t.Errorf("Algorithm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
