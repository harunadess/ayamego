package messaging

import (
	"github.com/bwmarrin/discordgo"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

// SendMessage sends a message to the channel provided
func SendMessage(session *discordgo.Session, channelID string, message string) {
	msg, err := session.ChannelMessageSend(channelID, message)
	if err != nil {
		logger.Error("SendMessage:", err)
		return
	}
	logger.Message(session.State.User.Username, "#", session.State.User.Discriminator, ": ", msg.Content)
}

// todo: sendMessageWithImage

// SendEmbedMessage sends an embedded message to the channel provided
func SendEmbedMessage(session *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) {
	msg, err := session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		logger.Error("SendEmbedMessage", err)
		return
	}
	logger.Message(session.State.User.Username, "#", session.State.User.Discriminator, ": ", msg.Content)
}
