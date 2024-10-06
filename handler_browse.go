package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CoDanTheBarbarian/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.input) > 2 {
		return fmt.Errorf("usage: browse <limit_integer>")
	} else if len(cmd.input) > 0 {
		n, err := strconv.Atoi(cmd.input[0])
		if err != nil {
			return fmt.Errorf("input must be a whole number")
		}
		limit = n
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed to get posts for user %s: %v", user.Name, err)
	}
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
