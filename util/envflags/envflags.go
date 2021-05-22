package envflags

import "flag"

// ConfigPath contains the file path to the main config for the bot
var ConfigPath *string

// BooruPath contains the file path to the booru config
var BooruPath *string

// DeviantPath contains the file path to the deviant config
var DeviantPath *string

func init() {
	ConfigPath = flag.String("config", "", "path to config .json file")
	BooruPath = flag.String("booru", "", "path to booru .json file")
	DeviantPath = flag.String("deviant", "", "path to deviant .json file")
	flag.Parse()
}
