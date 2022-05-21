package commands

import (
	"errors"
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
	"github.com/jordanjohnston/ayamego/util/timeparser"
)

// Prefix for standard commands
const Prefix string = "ayame, "
const HarunaUserId string = "150765618054823936"

const (
	Response_Complex  string = "COMPLEX"
	Response_Standard string = "STANDARD"
)

type commandFunction func(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string)

type commandHandler struct {
	description string
	exec        commandFunction
}

// todo: make the commandlers just have access to and send the messages directly from here
// will likely give more control over the response types and how responses are given

// todo: add other commands to bot
var commandlers = map[string]commandHandler{}

func init() {
	commandlers["help"] = commandHandler{
		description: "sends this message",
		exec:        generateHelpMessage,
	}
	commandlers["hello"] = commandHandler{
		description: "says hello to the user",
		exec:        sayHello,
	}
	commandlers["avatar"] = commandHandler{
		description: "returns your avatar image",
		exec:        avatar,
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
	commandlers["dice roll"] = commandHandler{
		description: "rolls dice up to the number specified",
		exec:        diceRoll,
	}
	commandlers["pick one"] = commandHandler{
		description: "pick one from supplied options. Format: <option1> | <option2> | <option3>",
		exec:        pickOne,
	}
	commandlers["add reminder"] = commandHandler{
		description: "sets a reminder for the specified time. Format: <reminder> @ DD/MM/YYYY HH:mm or <reminder> @ today HH:mm",
		exec:        addReminder,
	}
	commandlers["timestamp"] = commandHandler{
		description: "converts the specified date/time into a discord rich timestamp",
		exec:        createTimestampForDiscord,
	}
	commandlers["sleep"] = commandHandler{
		description: "asks bot go sleep",
		exec:        goSleep,
	}
}

// TryHandleStandardCommand checks if the message contains Prefix, and if it does
// tries to find and execute the appropriate command handler
func TryHandleStandardCommand(session *discordgo.Session, message *discordgo.MessageCreate) (bool, string, string) {
	response := "ayame does not know that command."
	responseType := Response_Standard

	content := strings.TrimPrefix(message.Content, Prefix)

	for k, f := range commandlers {
		if withoutPrefix := strings.TrimPrefix(content, k); withoutPrefix != content {
			response, responseType = f.exec(session, strings.Trim(withoutPrefix, " "), message)
			break
		}
	}
	return (response != ""), response, responseType
}

func sayHello(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	return "Hello " + discordMessage.Author.Mention() + "!", Response_Standard
}

func avatar(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	return discordMessage.Author.AvatarURL("256"), Response_Standard
}

func setActivity(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	err := discord.SetActivity(session, message)

	if err == nil {
		return "Successfully updated status!", Response_Standard
	}

	logger.Error("setActivity", err)
	return "There was an error updating status.", Response_Standard
}

// todo: on success this just returns an empty string because it's a command func, but this isn't great
// not sure how we refactor this out right now
// maybe instead of returning a string, we create a struct with a type + properties that contain the message to send
// for now, it just directly uses session to send a message
func booruSearch(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	found, results := booru.Search(message)
	logger.Info("Search output ", found, results)

	if !found {
		return "No results found for those search terms!", Response_Standard
	}

	embed := makeImageEmbed(results, "Powered by danbooru")
	_, err := session.ChannelMessageSendEmbed(discordMessage.ChannelID, embed)

	if err != nil {
		logger.Error("booruSearch", err)
		return "Command failed dazo, please check the logs", Response_Standard
	}
	// note: this does not log anything.. need to figure out how to do that from embed message
	logger.Message(session.State.User.Username, "#", session.State.User.Discriminator, ": ", results)

	return "", Response_Standard
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

func deviantSearch(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	found, results := deviant.Search(message)

	if !found {
		return "No results found for those search terms", Response_Standard
	}

	embed := makeImageEmbed(results, "Powered by deviantart")
	_, err := session.ChannelMessageSendEmbed(discordMessage.ChannelID, embed)

	if err != nil {
		logger.Error("deviantSearch", err)
		return "Command failed dazo, please check the logs", Response_Standard
	}
	// note: this does not log anything.. need to figure out how to do that from embed message
	logger.Message(session.State.User.Username, "#", session.State.User.Discriminator, ": ", results)

	return "", Response_Standard
}

func generateHelpMessage(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	response := "```html"
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

	return response, Response_Complex
}

func diceRoll(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	diceSides, err := strconv.Atoi(message)
	if err != nil {
		logger.Error("diceRoll", err)
		return "Command failed dazo, please check the command format and try again", Response_Standard
	}

	response := fmt.Sprintf("You rolled a %d!", rand.Intn(diceSides))

	return response, Response_Standard
}

func pickOne(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	options := strings.Split(message, "|")

	if len(options) < 2 {
		return "Command failed dazo, please check the command format and try again", Response_Standard
	}

	for i, v := range options {
		options[i] = strings.TrimSpace(v)
	}

	for _, v := range options {
		if len(v) < 1 {
			return "Command failed dazo, please check the command format and try again", Response_Standard
		}
	}

	go createDelayedResponse(options, session, discordMessage)

	return "docchi, docchi...", Response_Standard
}

func createDelayedResponse(options []string, session *discordgo.Session, discordMessage *discordgo.MessageCreate) {
	time.Sleep(time.Second * 1)
	delayedResponse := actuallyPickOne(options, discordMessage.Author.Mention())
	_, err := session.ChannelMessageSend(discordMessage.ChannelID, delayedResponse)
	if err != nil {
		logger.Error("failed to send delayed response", err)
	}
}

func actuallyPickOne(options []string, mentionStr string) string {
	rand.Seed(time.Now().UnixNano())
	randIdx := rand.Intn(len(options))

	return fmt.Sprintf("%v, Ayame has picked '%v' dayo", mentionStr, options[randIdx])
}

func addReminder(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	messageParts := make([]string, 2)

	for i, s := range strings.Split(message, "@") {
		messageParts[i] = strings.TrimSpace(s)
	}

	response := "Command failed dazo, please check the command format and try again"
	if len(messageParts) < 2 || (len(messageParts[0]) < 1 || len(messageParts[1]) < 1) {
		return response, Response_Standard
	}

	reminderText := messageParts[0]
	dateTimeStr := messageParts[1]

	timestampParts := strings.Split(dateTimeStr, " ")
	reminderTime, err := parseTimeOrTimestamp(timestampParts)

	if err != nil {
		logger.Error("error parsing time into format: ", err)
		return response, Response_Standard
	}

	response = fmt.Sprintf("Yo will remind you '%s' at %v!", reminderText, reminderTime.Format(timeparser.DateFormatExample))
	reminders.AddReminder(session, messageParts[0], discordMessage.Author.ID, time.Duration(reminderTime.Unix()))

	return response, Response_Standard
}

func parseTimeOrTimestamp(timeParts []string) (time.Time, error) {
	if len(timeParts) < 2 {
		return time.Time{}, errors.New("parseTimeFromCommand - not enough timeParts")
	}

	dateStr := timeParts[0]
	timeStr := timeParts[1]

	if strings.TrimSpace(dateStr) == "today" {
		parsedTime, err := timeparser.ParseTime(timeStr)
		if err != nil {
			logger.Error("error parsing time into format: ", err)
			return time.Time{}, err
		}
		return timeparser.SetTimeOfLeftToTimeOfRight(time.Now(), parsedTime), nil
	}

	parsedTime, err := timeparser.ParseDateTime(strings.Join(timeParts, " "))
	if err != nil {
		logger.Error("error parsing timestamp into format: ", err)
		return time.Time{}, nil
	}
	return parsedTime, nil
}

func createTimestampForDiscord(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	parts := strings.Split(message, " ")
	timestamp, err := parseTimeOrTimestamp(parts)
	if err != nil {
		logger.Error("error parsing time into format: ", err)
		return "Command failed dazo, please check the command format and try again", Response_Standard
	}
	return fmt.Sprintf("Here is the timestamp for <t:%v>: `<t:%v>`", timestamp.Unix(), timestamp.Unix()), Response_Standard
}

func goSleep(session *discordgo.Session, message string, discordMessage *discordgo.MessageCreate) (string, string) {
	response := "You do not have permission to do that"
	if discordMessage.Author.ID == HarunaUserId {
		response = "Otsunakiri!!"
		go doDisconnect(session)
	}

	return response, Response_Standard
}

func doDisconnect(session *discordgo.Session) {
	err := discord.Disconnect(session)
	if err != nil {
		logger.Fatal("error during disconnection", err)
	}
	logger.Info("exiting normally")
	os.Exit(0)
}
