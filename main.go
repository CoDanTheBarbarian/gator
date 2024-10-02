package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/CoDanTheBarbarian/gator/internal/config"
	"github.com/CoDanTheBarbarian/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	mainState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	mainCommands := commands{
		commandNameMap: make(map[string]func(*state, command) error),
	}
	mainCommands.register("login", handlerLogin)
	mainCommands.register("register", handlerRegister)
	mainCommands.register("reset", handlerReset)
	mainCommands.register("users", handlerListUsers)
	mainCommands.register("agg", handlerAggregate)

	if len(os.Args) < 2 {
		fmt.Println("invalid input")
		os.Exit(1)
	}
	commandName := os.Args[1]
	commandInput := []string{}
	if len(os.Args) > 2 {
		commandInput = os.Args[2:]
	}
	err = mainCommands.run(mainState, command{name: commandName, input: commandInput})
	if err != nil {
		log.Fatal(err)
	}
}
