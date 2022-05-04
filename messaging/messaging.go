package messaging

import (
	"github.com/bwmarrin/discordgo"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

// SendMessage sends a message to the channel provided
func SendMessage(session *discordgo.Session, channelID string, message string) {
	msg, err := session.ChannelMessageSend(channelID, message)
	if err != nil {
		logger.Error("SendMessage failed to send message: ", err)
		return
	}
	logger.Message(session.State.User.Username, "#", session.State.User.Discriminator, ": ", msg.Content)
}

// todo: sendMessageWithImage

// SendEmbedMessage sends an embedded message to the channel provided
func SendEmbedMessage(session *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) {
	msg, err := session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		logger.Error("SendEmbedMessage failed to send embed: ", err)
		return
	}
	logger.Message(session.State.User.Username, "#", session.State.User.Discriminator, ": ", msg.Content)
}

// SendMessageDM sends a message to the user provided
func SendMessageDM(session *discordgo.Session, userID string, message string) {
	channel, err := session.UserChannelCreate(userID)
	if err != nil {
		logger.Error("SendMessageDM, failed to create user DM: ", err)
		return
	}
	SendMessage(session, channel.ID, message)
}
