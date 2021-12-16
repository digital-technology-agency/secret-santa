package services

import "regexp"

const (
	//
	CmdExitGame      = "cmd_exit"
	CmdJoinGame      = "cmd_join"
	CmdLayerListGame = "cmd_player_list"
	CmdLanguageGame  = "cmd_language"
)

var (
	InitGameRegex, _ = regexp.Compile(".*дед.*мороз|санта|игр.*|хочу.*игр.*")
)
