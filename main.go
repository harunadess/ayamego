package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	discordactions "github.com/jordanjohnston/ayamego/discord/discordactions"
	discordinit "github.com/jordanjohnston/ayamego/discord/discordinit"
	"github.com/jordanjohnston/ayamego/messaging"
	envflags "github.com/jordanjohnston/ayamego/util/envflags"
	errors "github.com/jordanjohnston/ayamego/util/errors"
	logger "github.com/jordanjohnston/ayamego/util/logger"
)

var ayameSecrets struct {
	Token        string `json:"token"`
	DevChannelID string `json:"devChannelID"`
}

func init() {
	fPath := envflags.ConfigPath
	readConfig(fPath)
}

// func parseArgs() *string {
// 	configFilePath := flag.String("config", "", "path to config .json file")
// 	booruFilePath := flag.String("booru", "", "path to booru .json file")
// 	flag.Parse()

// 	if *configFilePath == "" {
// 		errors.FatalErrorHandler("parseArgs: ", fmt.Errorf("%v", "no -config specified"))
// 	}
// 	if *booruFilePath == "" {
// 		errors.FatalErrorHandler("parseArgs: ", fmt.Errorf("%v", "no -booru specified"))
// 	}

// 	return configFilePath
// }

func readConfig(fPath *string) {
	const maxJSONBytes int = 256

	file, err := os.Open(*fPath)
	defer file.Close()
	errors.FatalErrorHandler("readConfig: ", err)

	data := make([]byte, maxJSONBytes)
	count, err := file.Read(data)
	errors.FatalErrorHandler("readConfig: ", err)

	err = json.Unmarshal(data[:count], &ayameSecrets)
	errors.FatalErrorHandler("readConfig: ", err)
}

func main() {
	ayame := discordinit.SetupBot(ayameSecrets.Token)
	logger.Info("konnakiri!")

	messaging.SendMessage(ayame, ayameSecrets.DevChannelID, "yo dayo!!")

	err := discordactions.SetActivity(ayame, "playing Apex, probably..")
	errors.StandardErrorHandler("SetActvitiy: ", err)

	// wait here for ctrl+c or other signal end term
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close the discord session
	ayame.Close()
}
