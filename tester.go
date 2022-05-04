package main

import (
	"fmt"
	"time"
)

func main1() {

	dur1 := time.Duration(time.Now().Unix())
	c := make(chan bool, 1)

	time.AfterFunc(3*time.Second, func() {
		dur2 := time.Duration(time.Now().AddDate(0, 0, 1).Unix())
		diff := (dur2 - dur1) * time.Second

		fmt.Println("time diff", diff)
		c <- true
	})

	<-c
}
