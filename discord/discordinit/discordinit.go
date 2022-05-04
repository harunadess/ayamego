package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jordanjohnston/ayamego/commands"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

// SetupBot creates a discord session and opens the websocket
func SetupBot(token string) *discordgo.Session {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Fatal("error creating discord session:", err)
	}

	bot.Identify.Intents = discordgo.IntentsGuildMessages

	// attach handlers
	bot.AddHandler(commands.OnMessageCreate)

	err = bot.Open()
	if err != nil {
		logger.Fatal("error opening connection: ", err)
	}

	return bot
}
