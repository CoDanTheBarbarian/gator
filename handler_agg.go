package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/CoDanTheBarbarian/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	for i, item := range feedData.Channel.Item {
		pubDate, err := parseTimeString(feedData.Channel.Item[i].PubDate)
		if err != nil {
			log.Printf("Error parsing PubDate: %v", err)
			// handle the error or set a default value
			pubDate = time.Time{}
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: sql.NullTime{
				Time:  pubDate,
				Valid: err == nil,
			},
			FeedID: feed.ID,
		})
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok {
				// This is a PostgreSQL error
				if pgErr.Code == "23505" { // unique_violation
					// This is a duplicate URL, we can ignore this
					continue
				}
			}
			// This is some other type of error, we should log it
			log.Printf("Failed to create post for feed %s: %v", feed.Name, err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

func parseTimeString(pubDate string) (time.Time, error) {
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		"Mon, 02 Jan 2006 15:04:05 GMT",
		"Wed, 01 Sep 2022 00:00:00 GMT",
		"2006-01-02T15:04:05Z07:00",
	}

	for _, format := range formats {
		t, err := time.Parse(format, pubDate)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse time string: %s", pubDate)
}
