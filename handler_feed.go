package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/us0p/rss-aggregator/internal/database"
)

type parametersFeed struct {
    Name string `json:"name"`
    URL  string `json:"url"`
}

func (a *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
    decoder := json.NewDecoder(r.Body)
    params := parametersFeed{}
    err := decoder.Decode(&params)

    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }

    feed, err := a.DB.CreateFeed(r.Context(), database.CreateFeedParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: params.Name,
        Url: params.URL,
        UserID: user.ID,
    })

    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (a * apiConfig) handleGedFeeds(w http.ResponseWriter, r *http.Request) {
    feeds, err := a.DB.GetFeeds(r.Context())

    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Couldn't get list of feeds: %v", err))
        return
    }

    respondWithJSON(w, 200, databaseFeedsToFeeds(feeds))
}
