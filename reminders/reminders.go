package reminders

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jordanjohnston/ayamego/messaging"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

type Reminder struct {
	message    string
	remindTime time.Duration
	userId     string
}

// AddReminder adds a reminder to the current reminder list
func AddReminder(session *discordgo.Session, message string, userId string, remindTime time.Duration) {
	timeDiff := (remindTime - time.Duration(time.Now().Unix())) * time.Second

	reminder := Reminder{
		message:    message,
		userId:     userId,
		remindTime: timeDiff,
	}
	defer logger.Info("added reminder", reminder)

	time.AfterFunc(timeDiff, func() {
		triggerReminder(session, reminder)
	})
}

func triggerReminder(session *discordgo.Session, reminder Reminder) {
	logger.Info("triggered reminder: ", reminder.message, " ", reminder.userId, " ", reminder.remindTime)
	messaging.SendMessageDM(session, reminder.userId, reminder.message)
}
