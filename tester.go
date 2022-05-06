package main

import (
	"fmt"
	"time"
)

func test() {
	formatExample := "15:04"
	reminderTime := time.Now()

	parsedTime, err := time.ParseInLocation(formatExample, "20:18", time.Local)

	reminderTime = time.Date(reminderTime.Year(), reminderTime.Month(), reminderTime.Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)

	if err != nil {
		panic(err)
	}

	fmt.Println(reminderTime)
}

// func main() {
// 	test()
// }
