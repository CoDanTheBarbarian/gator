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

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.input) != 1 {
		return fmt.Errorf("usage: unfollow <feed_url>")
	}
	feed, err := s.db.GetFeedFromUrl(context.Background(), cmd.input[0])
	if err != nil {
		return err
	}
	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	followed := false
	for _, entries := range following {
		if entries.FeedID == feed.ID {
			followed = true
		}
	}
	if !followed {
		return fmt.Errorf("you are not currently following feed: %s", cmd.input[0])
	}
	err = s.db.DeleteFeedFollowEntry(context.Background(), database.DeleteFeedFollowEntryParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow %s error: %v", cmd.input[0], err)
	}
	fmt.Printf("You are no longer following feed: %v\n", feed.Name)
	return nil
}
