package main

import (
	"fmt"
	"time"
)

func main1() {

	loc, _ := time.LoadLocation("Europe/London")

	// This will look for the name CEST in the Europe/Berlin time zone.
	// const longForm = "Jan 2, 2006 3:04pm"
	const longForm = "02/01/2006 15:04"
	t, err := time.ParseInLocation(longForm, "05/06/2022 19:46", loc)
	fmt.Println(t)

	if err != nil {
		panic(err)
	}

	// Note: without explicit zone, returns time in given location.
	const shortForm = "2006-Jan-02"
	t, err = time.ParseInLocation(shortForm, "2012-Jul-09", loc)
	fmt.Println(t)

	if err != nil {
		panic(err)
	}

	fmt.Println(t)
}
