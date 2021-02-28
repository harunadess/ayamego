package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	discordactions "github.com/jordanjohnston/ayamego/discord/discordactions"
	discordinit "github.com/jordanjohnston/ayamego/discord/discordinit"
	"github.com/jordanjohnston/ayamego/messaging"
	errors "github.com/jordanjohnston/ayamego/util/errors"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

var ayameSecrets map[string]interface{}

func init() {
	fPath := parseArgs()
	readConfig(fPath)
}

func parseArgs() *string {
	fPath := flag.String("config", "", "path to config .json file")
	flag.Parse()

	if *fPath == "" {
		errors.FatalErrorHandler("parseArgs: ", fmt.Errorf("%v", "no -config specified"))
	}

	return fPath
}

func readConfig(fPath *string) {
	const maxJSONBytes int = 256

	file, err := os.Open(*fPath)
	errors.FatalErrorHandler("readConfig: ", err)

	data := make([]byte, maxJSONBytes)
	count, err := file.Read(data)
	errors.FatalErrorHandler("readConfig: ", err)

	err = json.Unmarshal(data[:count], &ayameSecrets)
	errors.FatalErrorHandler("readConfig: ", err)

	file.Close()
}

func main() {
	ayame := discordinit.SetupBot(ayameSecrets["token"].(string))
	logger.Info("konnakiri!")

	messaging.SendMessage(ayame, ayameSecrets["generalID"].(string), "yo dayo!!")

	err := discordactions.SetActivity(ayame, "playing", "Apex, probably..")
	errors.StandardErrorHandler("SetActvitiy: ", err)

	// wait here for ctrl+c or other signal end term
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close the discord session
	ayame.Close()
}
