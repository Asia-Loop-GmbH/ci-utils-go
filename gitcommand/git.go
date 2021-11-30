package gitcommand

import (
	"log"
	"os"
)

func GitCommand() {
	subCommand := os.Args[2]

	switch subCommand {
	case "create-release":
		createRelease()
	default:
		log.Panicf("unsupported subcommand '%s' for 'git'", subCommand)
	}
}
