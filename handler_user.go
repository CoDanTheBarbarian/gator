package main

import (
	"context"
	"fmt"
	"time"

	"github.com/CoDanTheBarbarian/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.input) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	name := cmd.input[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
	}
	fmt.Printf("username: %s has been set.\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.input) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.input[0],
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set created user, error: %w", err)
	}
	fmt.Println("User successfully created")
	printUser(user)
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	items, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not get users from database: %w", err)
	}
	for _, user := range items {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%v (current)\n", user.Name)
			continue
		}
		fmt.Printf("%v\n", user.Name)
	}
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
