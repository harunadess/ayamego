package discord

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jordanjohnston/harunago/commands"
	logger "github.com/jordanjohnston/harunago/util"
)

// SetupBot creates a discord session and opens the websocket
func SetupBot(token string) *discordgo.Session {
	bot, err := discordgo.New("Bot " + token)
	errorHandler("error creating discord session:", err)

	bot.Identify.Intents = discordgo.IntentsGuildMessages

	// attach handlers
	bot.AddHandler(commands.OnMessageCreate)

	err = bot.Open()
	errorHandler("error opening connection: ", err)

	return bot
}

func errorHandler(msg string, err error) {
	if err != nil {
		logger.Error(msg, err)
		os.Exit(1)
	}
}

// SetActivity sets the activity of the bot
func SetActivity(session *discordgo.Session, activityType string, activityMsg string) error {
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
