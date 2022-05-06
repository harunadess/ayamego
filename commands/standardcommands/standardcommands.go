package commands

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jordanjohnston/ayamego/booru"
	"github.com/jordanjohnston/ayamego/deviant"
	discord "github.com/jordanjohnston/ayamego/discord/discordactions"
	"github.com/jordanjohnston/ayamego/imageresults"
	"github.com/jordanjohnston/ayamego/reminders"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

// Prefix for standard commands
const Prefix string = "ayame, "
const HarunaUserId string = "150765618054823936"

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
		description: "( ͡° ͜ʖ ͡°)",
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
	commandlers["add reminder"] = commandHandler{
		description: "sets a reminder for the specified time. Format: <reminder> @ DD/MM/YYYY HH:mm",
		exec:        addReminder,
	}
	commandlers["sleep"] = commandHandler{
		description: "asks bot go sleep",
		exec:        goSleep,
	}
}

// TryHandleStandardCommand checks if the message contains Prefix, and if it does
// tries to find and execute the appropriate command handler
func TryHandleStandardCommand(session *discordgo.Session, message *discordgo.MessageCreate) (bool, string) {
	response := "ayame does not know that command."

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

	logger.Error("setActivity", err)
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
		logger.Error("booruSearch", err)
		return "Command failed dazo, please check the logs"
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
		logger.Error("deviantSearch", err)
		return "Command failed dazo, please check the logs"
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
		if k != "deviant" {
			response += fmt.Sprintf("\n<%s>: %s\n", k, v.description)
		} else if discordMessage.Author.ID == HarunaUserId {
			response += fmt.Sprintf("\n<%s>: %s\n", k, v.description)
		}
	}
	response += "\n===============================\n```"

	return response
}

func diceRoll(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	diceSides, err := strconv.Atoi(message)
	if err != nil {
		logger.Error("diceRoll", err)
		return "Command failed dazo, please check the command format and try again"
	}

	response := fmt.Sprintf("You rolled a %d!", rand.Intn(diceSides))

	return response
}

func addReminder(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	messageParts := make([]string, 2)

	for i, s := range strings.Split(message, "@") {
		messageParts[i] = strings.TrimSpace(s)
	}

	response := "Command failed dazo, please check the command format and try again"
	if len(messageParts) < 2 || (len(messageParts[0]) < 1 || len(messageParts[1]) < 1) {
		return response
	}

	reminderText := messageParts[0]
	timeStr := messageParts[1]

	// note: this time format example must be 2nd Jan 2006 @ 15:04
	dateFormatExample := "02/01/2006 15:04"
	timeFormatExample := "15:04"
	reminderTime := time.Now()

	dateText := strings.Split(timeStr, " ")
	if strings.TrimSpace(dateText[0]) == "today" {
		parsedTime, err := time.ParseInLocation(timeFormatExample, strings.TrimSpace(timeStr), time.Local)
		if err != nil {
			logger.Error("error parsing time into format: ", err)
			return response
		}
		reminderTime = setClockFromSecondTime(reminderTime, parsedTime)
	} else {
		parsedTime, err := time.ParseInLocation(dateFormatExample, strings.TrimSpace(timeStr), time.Local)
		if err != nil {
			logger.Error("error parsing time into format: ", err)
			return response
		}
		reminderTime = parsedTime
	}

	response = fmt.Sprintf("Yo will remind you to '%s' at %v!", reminderText, reminderTime.Format(dateFormatExample))
	reminders.AddReminder(session, messageParts[0], discordMessage.Author.ID, time.Duration(reminderTime.Unix()))

	return response
}

func setClockFromSecondTime(t1 time.Time, t2 time.Time) time.Time {
	return time.Date(t1.Year(), t1.Month(), t1.Day(),
		t2.Hour(), t2.Minute(), 0, 0, time.Local)
}

func goSleep(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) string {
	response := "You do not have permission to do that"
	if discordMessage.Author.ID == HarunaUserId {
		response = "Otsunakiri"
		defer doDisconnect(session)
	}

	return response
}

func doDisconnect(session *discordgo.Session) {
	err := discord.Disconnect(session)
	if err != nil {
		logger.Fatal("error during disconnection", err)
	}
	logger.Info("exiting normally")
	os.Exit(0)
}
