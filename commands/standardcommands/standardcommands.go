package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	discord "github.com/jordanjohnston/ayamego/discord/discordactions"
	errors "github.com/jordanjohnston/ayamego/util/errors"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

// Prefix for standard commands
const Prefix string = "ayame, "

type commandFunction func(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string

type commandHandler struct {
	description string
	exec        commandFunction
}

// todo: add other commands to bot
var commandlers = map[string]commandHandler{}

func init() {
	commandlers["hello"] = commandHandler{
		description: "says hello to the user",
		exec:        sayHello,
	}
	commandlers["setActivity"] = commandHandler{
		description: "sets bot activity",
		exec:        setActivity,
	}
	commandlers["search for"] = commandHandler{
		description: "search danbooru for an image",
		exec:        booruSearch,
	}
	commandlers["help"] = commandHandler{
		description: "sends this message",
		exec:        generateHelpMessage,
	}
	commandlers["dice roll"] = commandHandler{
		description: "rolls dice up to the number specified",
		exec:        diceRoll,
	}
}

// TryHandleStandardCommand checks if the message contains Prefix, and if it does
// tries to find and execute the appropriate command handler
func TryHandleStandardCommand(session *discordgo.Session, message *discordgo.MessageCreate) (bool, string) {
	response := ""

	content := strings.TrimPrefix(message.Content, Prefix)

	for k, f := range commandlers {
		if withoutPrefix := strings.TrimPrefix(content, k); withoutPrefix != content {
			response = f.exec(session, strings.Trim(withoutPrefix, " "), message)
			break
		}
	}
	return (response != ""), response
}

func sayHello(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	return "Hello " + discordMessage.Author.Mention() + "!"
}

func setActivity(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	err := discord.SetActivity(session, message)

	if err == nil {
		return "Successfully updated status!"
	}

	logger.Error(err)
	return "There was an error updating status."
}

func booruSearch(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	// results :=
	return ""
}

func generateHelpMessage(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	response := "```css"
	response += "\n======== Help Commands ========"
	response += "\nPrefix any command with 'ayame,'\n"

	for k, v := range commandlers {
		response += fmt.Sprintf("\n<%s>: %s\n", k, v.description)
	}
	response += "\n===============================\n```"

	return response
}

func diceRoll(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	diceSides, err := strconv.Atoi(message)
	errors.StandardErrorHandler("diceRoll", err)

	response := fmt.Sprintf("You rolled a %d!", rand.Intn(diceSides))

	return response
}
