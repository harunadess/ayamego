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
	hasResponse, response, responseType := false, "", standardcommands.Response_Standard

	if strings.HasPrefix(message.Content, standardcommands.Prefix) {
		// do standard command things
		hasResponse, response, responseType = standardcommands.TryHandleStandardCommand(session, message)
	} else if strings.HasPrefix(message.Content, musicCommandPrefix) {
		hasResponse, response = false, "music command"
	} else {
		hasResponse, response = substringcommands.TryHandleSubstringCommand(message)
		if hasResponse {
			messaging.Reply(session, channelID, message.Reference(), response)
			return
		}
	}

	if hasResponse {
		logger.Message(message.Author, ": ", message.Content)
		switch responseType {
		case standardcommands.Response_Standard:
			messaging.SendMessage(session, channelID, response)
		case standardcommands.Response_Complex:
			messaging.SendMessageComplex(session, channelID, response)
		}
	}
}
