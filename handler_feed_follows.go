package main

import (
	"context"
	"fmt"
	"time"

	"github.com/CoDanTheBarbarian/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.input) != 1 {
		return fmt.Errorf("usage: follow <url>")
	}
	feed, err := s.db.GetFeedFromUrl(context.Background(), cmd.input[0])
	if err != nil {
		return fmt.Errorf("failed to get feed from url: %v - add feed to database", cmd.input[0])
	}
	row, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed_follow entry: %v", err)
	}
	fmt.Printf("Feed: %v followed by user %v\n", row.FeedName, row.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.input) > 0 {
		return fmt.Errorf("no argument required for following command")
	}
	followRows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feeds for user %s: %v", user.Name, err)
	}
	for _, row := range followRows {
		fmt.Printf("Feed: %v - User: %v\n", row.FeedName, row.UserName)
	}
	return nil
}
