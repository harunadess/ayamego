package reminders

import (
	"testing"
	"time"

	"github.com/jordanjohnston/ayamego/reminders"
)

func TestAddReminder(t *testing.T) {
	reminders.AddReminder(nil, "test", "1", time.Duration(time.Now().Unix()))
}
