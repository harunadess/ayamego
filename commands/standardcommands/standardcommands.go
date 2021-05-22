package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jordanjohnston/ayamego/booru"
	"github.com/jordanjohnston/ayamego/deviant"
	discord "github.com/jordanjohnston/ayamego/discord/discordactions"
	"github.com/jordanjohnston/ayamego/imageresults"
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
	commandlers["deviant"] = commandHandler{
		description: "search deviantart for an image",
		exec:        deviantSearch,
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

// todo: on success this just returns an empty string because it's a command func, but this isn't great
// not sure how we refactor this out right now
// maybe instead of returning a string, we create a struct with a type + properties that contain the message to send
// for now, it just directly uses session to send a message
func booruSearch(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	found, results := booru.Search(message)
	logger.Info("Search output ", found, results)

	if !found {
		return "No results found for those search terms!"
	}

	embed := makeImageEmbed(results, "Powered by danbooru")
	_, err := session.ChannelMessageSendEmbed(discordMessage.ChannelID, embed)

	if err != nil {
		logger.Error("messaging: ", err)
	}
	// note: this does not log anything.. need to figure out how to do that from embed message
	logger.Message(session.State.User.Username, "#", session.State.User.Discriminator, ": ", results)

	return ""
}

func makeImageEmbed(results imageresults.SearchResults, footerText string) *discordgo.MessageEmbed {
	const msgColor int = 16750848

	msg := discordgo.MessageEmbed{
		URL:         results.Images.ImageURL,
		Type:        discordgo.EmbedTypeImage,
		Title:       results.Title,
		Description: results.Tags,
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       msgColor,
		Footer:      &discordgo.MessageEmbedFooter{Text: footerText},
		Image:       &discordgo.MessageEmbedImage{URL: results.Images.Thumbnail, Height: 720, Width: 576},
	}

	return &msg
}

func deviantSearch(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	found, results := deviant.Search(message)

	if !found {
		return "No results found for those search terms"
	}

	embed := makeImageEmbed(results, "Powered by deviantart")
	_, err := session.ChannelMessageSendEmbed(discordMessage.ChannelID, embed)

	if err != nil {
		logger.Error("messaging: ", err)
	}
	// note: this does not log anything.. need to figure out how to do that from embed message
	logger.Message(session.State.User.Username, "#", session.State.User.Discriminator, ": ", results)

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
