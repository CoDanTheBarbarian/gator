package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CoDanTheBarbarian/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}
	mainState := state{}
	mainState.cfg = &cfg

	mainCommands := commands{
		commandNameMap: make(map[string]func(*state, command) error),
	}
	mainCommands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("invalid input")
		os.Exit(1)
	}
	commandName := os.Args[1]
	commandInput := []string{}
	if len(os.Args) > 2 {
		commandInput = os.Args[2:]
	}
	err = mainCommands.run(&mainState, command{name: commandName, input: commandInput})
	if err != nil {
		log.Fatal(err)
	}
}
