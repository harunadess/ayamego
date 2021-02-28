package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jordanjohnston/ayamego/commands"
	errors "github.com/jordanjohnston/ayamego/util/errors"
)

// SetupBot creates a discord session and opens the websocket
func SetupBot(token string) *discordgo.Session {
	bot, err := discordgo.New("Bot " + token)
	errors.FatalErrorHandler("error creating discord session:", err)

	bot.Identify.Intents = discordgo.IntentsGuildMessages

	// attach handlers
	bot.AddHandler(commands.OnMessageCreate)

	err = bot.Open()
	errors.FatalErrorHandler("error opening connection: ", err)

	return bot
}
