package logger

import (
	"fmt"
	"time"
)

var tags = map[string]string{
	"info":    "INFO",
	"error":   "ERR ",
	"message": "MSG ",
	"image":   "IMG ",
}

// Info (a ...interface): logs with INFO tag
func Info(a ...interface{}) {
	log(tags["info"], a...)
}

// Error (a ...interface): logs with ERR tag
func Error(a ...interface{}) {
	log(tags["error"], a...)
}

// Message (a ...interface): logs with MSG tag
func Message(a ...interface{}) {
	log(tags["message"], a...)
}

// Image (a ...interface): logs with IMG tag
func Image(a ...interface{}) {
	log(tags["image"], a...)
}

func log(tag string, a ...interface{}) {
	prefix := fmt.Sprintf("%v [%v]", time.Now().Format(time.Stamp), tag)
	args := fmt.Sprint(a...)
	fmt.Println(prefix, args)
}
