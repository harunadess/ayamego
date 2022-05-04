package util

import (
	"fmt"
	"os"
	"time"
)

var tags = map[string]string{
	"info":    "INFO",
	"error":   "ERR ",
	"message": "MSG ",
	"image":   "IMG ",
	"fatal":   "FATL",
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

// Fatal (a ...interface): logs the error with the FATL tag, before immediately exiting the program with error code 1
func Fatal(a ...interface{}) {
	log(tags["fatal"], a...)
	os.Exit(1)
}

func log(tag string, a ...interface{}) {
	prefix := fmt.Sprintf("%v [%v]", time.Now().Format(time.Stamp), tag)
	args := fmt.Sprint(a...)
	fmt.Println(prefix, args)
}
