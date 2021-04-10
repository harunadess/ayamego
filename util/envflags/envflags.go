package envflags

import "flag"

// ConfigPath contains the file path to the main config for the bot
var ConfigPath *string

// BooruPath contains the file path to the booru config
var BooruPath *string

func init() {
	ConfigPath = flag.String("config", "", "path to config .json file")
	BooruPath = flag.String("booru", "", "path to booru .json file")
	flag.Parse()
}
