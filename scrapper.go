package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/us0p/rss-aggregator/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
    log.Printf("Scraping on %d goroutines every %s duration", concurrency, timeBetweenRequest)
    ticker := time.NewTicker(timeBetweenRequest)
    // Here we are passing only the increment part of the for loop to cause the loop to execute
    // immediatly on the first call, and then wait for the next tick.
    for ; ; <-ticker.C {
        feeds, err := db.GetNextFeedToFetch(
            context.Background(),
            int32(concurrency),
        )

        if err != nil {
            log.Println("error fetching feeds:", err)
            continue
        }

        wg := &sync.WaitGroup{}
        for _, feed := range feeds {
            wg.Add(1)

            go scrapFeed(wg, db, feed)
        }

        wg.Wait()
    }
}

func scrapFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed){
    defer wg.Done()

    _, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
    if err != nil {
        log.Println("error marking feed as fetched:", err)
        return 
    }

    rssFeed, err := urlToFeed(feed.Url)
    if err != nil {
        log.Println("error fetching feed:", err)
        return 
    }

    for _, item := range rssFeed.Channel.Item {
        description := sql.NullString{}

        if item.Description != "" {
            description.String = item.Description
            description.Valid = true
        }

        t, err := time.Parse(time.RFC1123Z, item.PubDate)

        if err != nil {
            log.Printf("Error parsing date: %s, %v\n", item.PubDate, err)
            continue
        }

        _, errPost := db.CreatePost(context.Background(), database.CreatePostParams{
            ID: uuid.New(),
            CreatedAt: time.Now().UTC(),
            UpdatedAt: time.Now().UTC(),
            FeedID: feed.ID,
            Title: item.Title,
            Description: description,
            PublishedAt: t,
            Url: item.Link,
        })

        if errPost != nil {
            if strings.Contains(errPost.Error(), "duplicate key"){
                continue
            }
            log.Println("error creating post:", errPost)
        }
    }

    log.Printf("Feed %s collected %d posts found\n", feed.Name, len(rssFeed.Channel.Item))
}
