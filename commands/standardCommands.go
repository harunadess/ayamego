package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	logger "github.com/jordanjohnston/harunago/util"
)

const commandPrefix string = "ayame, "

// OnMessageCreate is a handler for the discord event MessageCreate
func OnMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	channelID := message.ChannelID
	response := ""

	response = specificMessageHandler(message.Content)
	if response != "" {
		sendMessage(session, channelID, response)
		return
	}

	if strings.HasPrefix(message.Content, commandPrefix) {
		logger.Message(message)
		response = basicCommandHandler(session, message)
	}

	if response != "" {
		sendMessage(session, channelID, response)
	}
}

func specificMessageHandler(content string) string {
	response := ""

	switch content {
	case "yo":
		response = "dayo!"
	case "konnakiri":
		response = "konnakiri!"
	case "hold on":
		response = "chotto machete!"
	}

	return response
}

func sendMessage(session *discordgo.Session, channelID string, message string) {
	msg, err := session.ChannelMessageSend(channelID, message)
	handleGenericError(err)
	logger.Message(msg.Content)
}

func basicCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) string {
	msgContent := message.Content[len(commandPrefix):]
	response := ""

	if strings.Contains(msgContent, "setActivity") {
		content := msgContent[len("setActivity "):]
		activityType := determineActivity(content)
		status := content[len(activityType)+1:]
		err := setActivity(session, activityType, status)
		handleGenericError(err)
		response = "Successfully updated status!"
	}

	return response
}

func determineActivity(content string) string {
	activity := ""

	if strings.Contains(content, "playing") {
		activity = "playing"
	} else if strings.Contains(content, "listening") {
		activity = "listening"
	} else if strings.Contains(content, "idle") {
		activity = "idle"
	}

	return activity
}

func setActivity(session *discordgo.Session, activityType string, activityMsg string) error {
	idle := 0
	var err error

	lcActivityType := strings.ToLower(activityType)

	switch lcActivityType {
	case "idle":
		idle = 1
		err = session.UpdateGameStatus(idle, activityMsg)
	case "playing":
		err = session.UpdateGameStatus(idle, activityMsg)
	case "listening":
		err = session.UpdateListeningStatus(activityMsg)
	}

	return err
}

func handleGenericError(err error) {
	if err != nil {
		logger.Error("standardCommands: ", err)
	}
}
