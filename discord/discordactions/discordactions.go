package discord

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

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
	default:
		err = errors.New("Unrecognised activity type: " + activityType)
	}

	standardErrorHandler("SetActivity: ", err)

	if err == nil {
		logger.Info("Set activity to: ", activityMsg)
	}

	return err
}

func standardErrorHandler(msg string, err error) {
	if err != nil {
		logger.Error(msg, err)
	}
}
