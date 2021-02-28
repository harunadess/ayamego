package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	discord "github.com/jordanjohnston/ayamego/discord/discordactions"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

// Prefix for standard commands
const Prefix string = "ayame, "

type commandHandler func(session *discordgo.Session, message string) string

// todo: add other commands to bot
var commandlers = map[string]commandHandler{
	"setActivity ": setActivity,
	"search for ":  booruSearch,
}

// TryHandleStandardCommand checks if the message contains Prefix, and if it does
// tries to find and execute the appropriate command handler
func TryHandleStandardCommand(session *discordgo.Session, message *discordgo.MessageCreate) (bool, string) {
	response := ""

	content := strings.TrimPrefix(message.Content, Prefix)

	for k, f := range commandlers {
		if withoutPrefix := strings.TrimPrefix(content, k); withoutPrefix != content {
			response = f(session, withoutPrefix)
			break
		}
	}
	return (response != ""), response
}

func setActivity(session *discordgo.Session, message string) string {
	err := discord.SetActivity(session, message)

	if err == nil {
		return "Successfully updated status!"
	} else {
		logger.Error(err)
		return "There was an error updating status."
	}
}

func booruSearch(session *discordgo.Session, message string) string {
	return ""
}
