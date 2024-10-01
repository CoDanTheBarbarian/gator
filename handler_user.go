package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.input) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	err := s.cfg.SetUser(cmd.input[0])
	if err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
	}
	fmt.Printf("username: %s has been set.\n", cmd.input[0])
	return nil
}
