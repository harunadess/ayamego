package discord

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type activity struct {
	idle         int
	activityType string
	msg          string
}

var activityTypes = []string{
	"idle",
	"playing",
	"listening",
	"streaming",
}

const (
	idleActivity = iota
	playingActivity
	listeningActivity
	streamingActivity
)

const streamingURL = "https://twitch.tv/harunadess"

// SetActivity sets the activity of the bot
func SetActivity(session *discordgo.Session, message string) error {

	var err error

	updatedActivity := activity{}

	for _, v := range activityTypes {
		if withoutPrefix := strings.TrimPrefix(message, v); withoutPrefix != message {
			updatedActivity.activityType = v
			if v == activityTypes[0] {
				updatedActivity.idle = 999
			}
			updatedActivity.msg = strings.Trim(withoutPrefix, " ")
		}
	}

	switch updatedActivity.activityType {
	case activityTypes[idleActivity]:
		fallthrough
	case activityTypes[playingActivity]:
		err = session.UpdateGameStatus(updatedActivity.idle, updatedActivity.msg)
	case activityTypes[listeningActivity]:
		err = session.UpdateListeningStatus(updatedActivity.msg)
	case activityTypes[streamingActivity]:
		err = session.UpdateStreamingStatus(updatedActivity.idle, updatedActivity.msg, streamingURL)
	default:
		err = errors.New("unrecognised activity type given:" + updatedActivity.activityType)
	}

	return err
}

func Disconnect(session *discordgo.Session) error {
	session.Client.CloseIdleConnections()
	return session.Close()
}
