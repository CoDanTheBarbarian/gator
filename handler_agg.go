package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CoDanTheBarbarian/gator/internal/database"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.input) < 1 {
		return fmt.Errorf("usage: agg <duration_string>")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.input[0])
	if err != nil {
		return fmt.Errorf("input is not a duration string: %v", err)
	}
	fmt.Printf("Collecting feeds every %s", cmd.input[0])
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Couldn't get next feed to fetch: %v", err)
	}
	log.Printf("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
	}
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't fetch feed %s: %v", feed.Name, err)
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
