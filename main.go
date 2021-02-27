package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jordanjohnston/harunago/discord"
	logger "github.com/jordanjohnston/harunago/util"
)

var ayameSecrets map[string]interface{}

func init() {
	fPath := parseArgs()
	readConfig(fPath)
}

// parseArgs reads the command line flags specified
func parseArgs() *string {
	fPath := flag.String("config", "", "path to config .json file")
	flag.Parse()

	if *fPath == "" {
		handleGenericError(fmt.Errorf("%v", "no -config specified"))
	}

	return fPath
}

func handleGenericError(err error) {
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

// readConfig reads the json file from the -config command line
// arg specified
func readConfig(fPath *string) {
	const maxJSONBytes int = 256

	file, err := os.Open(*fPath)
	handleGenericError(err)

	data := make([]byte, maxJSONBytes)
	count, err := file.Read(data)
	handleGenericError(err)

	err = json.Unmarshal(data[:count], &ayameSecrets)
	handleGenericError(err)

	file.Close()
}

// todo: split things out into individual packages
// i.e. move discord specific stuff into it's own file
// move handlers into specific packages too
func main() {
	ayame := discord.SetupBot(ayameSecrets["token"].(string))
	logger.Info("konnakiri!")

	msg, _ := ayame.ChannelMessageSend(ayameSecrets["generalID"].(string), "yo dayo!!")
	logger.Message(msg.Content)

	err := discord.SetActivity(ayame, "playing", "Apex probably")
	if err != nil {
		logger.Error(err)
	}

	// wait here for ctrl+c or other signal end term
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close the discord session
	ayame.Close()
}
