package main

import (
	"asialoop.de/ci-utils-go/gitcommand"
	"flag"
	"log"
	"os"
)

func main() {
	flag.Parse()
	command := os.Args[1]

	switch command {
	case "git":
		gitcommand.GitCommand()
	default:
		log.Panicf("unsupported command '%s'", command)
	}
}
