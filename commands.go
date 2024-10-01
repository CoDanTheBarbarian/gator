package main

import (
	"fmt"
)

type command struct {
	name  string
	input []string
}

type commands struct {
	commandNameMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if _, exists := c.commandNameMap[name]; exists {
		fmt.Printf("command '%s' is already registered\n", name)
		return
	}
	c.commandNameMap[name] = f
	fmt.Printf("command '%s' successfully registered\n", name)
}

func (c *commands) run(s *state, cmd command) error {
	commandFunction, ok := c.commandNameMap[cmd.name]
	if !ok {
		return fmt.Errorf("command %v does not exist", cmd.name)
	}
	return commandFunction(s, cmd)
}
