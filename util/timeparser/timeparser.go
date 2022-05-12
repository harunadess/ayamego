package timeparser

import (
	"errors"
	"strings"
	"time"
)

// note: this time format example must be 2nd Jan 2006 @ 15:04
const DateFormatExample string = "02/01/2006 15:04"
const TimeFormatExample string = "15:04"
const TimestampFormatExample string = "02/01/2006 15:04"

func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "today" {
		return time.Now(), nil
	}
	if len(dateStr) < len(DateFormatExample) {
		return time.Time{}, errors.New("timeparser - failed to parse date string")
	}
	return time.ParseInLocation(DateFormatExample, strings.TrimSpace(dateStr), time.Local)
}

func ParseTime(timeStr string) (time.Time, error) {
	if timeStr == "today" {
		return time.Now(), nil
	}
	if len(timeStr) < len(TimeFormatExample) {
		return time.Time{}, errors.New("timeparser - failed to parse time string")
	}
	return time.ParseInLocation(TimeFormatExample, strings.TrimSpace(timeStr), time.Local)
}

func ParseDateTime(timestamp string) (time.Time, error) {
	if timestamp == "today" {
		return time.Now(), nil
	}
	if len(timestamp) < len(TimestampFormatExample) {
		return time.Time{}, errors.New("timeparser - failed to parse timestamp")
	}
	return time.ParseInLocation(TimestampFormatExample, strings.TrimSpace(timestamp), time.Local)
}

func SetTimeOfLeftToTimeOfRight(t1 time.Time, t2 time.Time) time.Time {
	return time.Date(t1.Year(), t1.Month(), t1.Day(),
		t2.Hour(), t2.Minute(), 0, 0, time.Local)
}
