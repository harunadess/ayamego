package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Prefix for standard commands
const Prefix string = "ayame, "

type commandHandler func(session *discordgo.Session, message string) string

var commandlers = map[string]commandHandler{
	"setActivity": setActivity,
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
	return "Successfully updated activity!"
}

/*
	Todo: need to figure out a more elegant solution for this
*/
// func basicCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) string {
// 	msgContent := message.Content[len(Prefix):]
// 	response := ""

// 	if strings.Contains(msgContent, "setActivity") {
// 		content := msgContent[len("setActivity "):]
// 		activityType := determineActivity(content)
// 		status := content[len(activityType)+1:]
// 		err := setActivity(session, activityType, status)
// 		handleGenericError(err)
// 		response = "Successfully updated status!"
// 	}

// 	return response
// }

// func determineActivity(content string) string {
// 	activity := ""

// 	if strings.Contains(content, "playing") {
// 		activity = "playing"
// 	} else if strings.Contains(content, "listening") {
// 		activity = "listening"
// 	} else if strings.Contains(content, "idle") {
// 		activity = "idle"
// 	}

// 	return activity
// }

// func setActivity(session *discordgo.Session, activityType string, activityMsg string) error {
// 	idle := 0
// 	var err error

// 	lcActivityType := strings.ToLower(activityType)

// 	switch lcActivityType {
// 	case "idle":
// 		idle = 1
// 		err = session.UpdateGameStatus(idle, activityMsg)
// 	case "playing":
// 		err = session.UpdateGameStatus(idle, activityMsg)
// 	case "listening":
// 		err = session.UpdateListeningStatus(activityMsg)
// 	}

// 	return err
// }
