package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var wholeMessages = map[string]string{
	"hold on":    "chotto machete!",
	"konnakiri":  "konnakiri!",
	"otsunakiri": "otsunakiri deshita!",
}

var substrings = map[string]string{
	"yo": "dayo!",
}

// TryHandleSubstringCommand tries to parse a substring command
// returns a bool which will be false if no response was generated
func TryHandleSubstringCommand(message *discordgo.MessageCreate) (bool, string) {
	content := message.Content

	response := findSubstringFromMap(content, wholeMessages)
	if response == "" {
		response = findSubstringFromMap(content, substrings)
	}

	return (response != ""), response
}

func findSubstringFromMap(str string, m map[string]string) string {
	for k, v := range m {
		if strings.Contains(str, k) {
			return v
		}
	}
	return ""
}
