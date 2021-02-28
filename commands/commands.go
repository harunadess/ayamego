package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	standardcommands "github.com/jordanjohnston/ayamego/commands/standardcommands"
	substringcommands "github.com/jordanjohnston/ayamego/commands/substringcommands"
	"github.com/jordanjohnston/ayamego/messaging"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

const musicCommandPrefix string = "+"

// OnMessageCreate is a handler for the discord event MessageCreate
func OnMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	channelID := message.ChannelID
	hasResponse, response := false, ""

	if strings.HasPrefix(message.Content, standardcommands.Prefix) {
		// do standard command things
		hasResponse, response = standardcommands.TryHandleStandardCommand(session, message)
	} else if strings.HasPrefix(message.Content, musicCommandPrefix) {
		hasResponse, response = false, "music command"
	} else {
		hasResponse, response = substringcommands.TryHandleSubstringCommand(message)
	}

	if hasResponse {
		logger.Message(message.Author, ": ", message.Content)
		messaging.SendMessage(session, channelID, response)
	}
}
