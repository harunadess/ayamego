package messaging

import (
	"github.com/bwmarrin/discordgo"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

// SendMessage sends a message to the channel provided
func SendMessage(session *discordgo.Session, channelID string, message string) {
	msg, err := session.ChannelMessageSend(channelID, message)
	handleGenericError(err)
	logger.Message(session.State.User.Username, "#", session.State.User.Discriminator, ": ", msg.Content)
}

// todo: sendMessageWithImage
// todo: sendMessageWithEmbed

func handleGenericError(err error) {
	if err != nil {
		logger.Error("messaging: ", err)
	}
}
