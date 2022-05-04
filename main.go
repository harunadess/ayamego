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

func readConfig(fPath *string) {
	const maxJSONBytes int = 256

	file, err := os.Open(*fPath)
	if err != nil {
		logger.Fatal("readConfig: ", err)
	}
	defer file.Close()

	data := make([]byte, maxJSONBytes)
	count, err := file.Read(data)
	if err != nil {
		logger.Fatal("readConfig: ", err)
	}

	err = json.Unmarshal(data[:count], &ayameSecrets)
	if err != nil {
		logger.Fatal("readConfig: ", err)
	}
}

func main() {
	ayame := discordinit.SetupBot(ayameSecrets.Token)
	logger.Info("konnakiri!")

	messaging.SendMessage(ayame, ayameSecrets.DevChannelID, "yo dayo!!")

	err := discordactions.SetActivity(ayame, "playing Apex, probably..")
	if err != nil {
		logger.Error("SetActvitiy: ", err)
	}

	// wait here for ctrl+c or other signal end term
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close the discord session
	ayame.Close()
}
