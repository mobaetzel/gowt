package main

import (
	"log"
	"os"
	"word_template_service/src/commands"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing command. commands are: serve, process")
	}

	cmd := os.Args[1]

	switch cmd {
	case "serve":
		commands.Serve()
	case "process":
		commands.Process()
	default:
		log.Fatalf("invalid command %s. commands are: serve, process", cmd)
	}
}
