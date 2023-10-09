package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"
	"github.com/us0p/rss-aggregator/internal/database"
)

type parametersFeedFollows struct {
    FeedID  uuid.UUID `json:"feed_id"`
}

func (a *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
    decoder := json.NewDecoder(r.Body)
    params := parametersFeedFollows{}
    err := decoder.Decode(&params)

    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }

    feedFollow, err := a.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        UserID: user.ID,
        FeedID: params.FeedID,
    })

    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (a *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
    feedFollows, err := a.DB.GetFeedFollows(r.Context(), user.ID)

    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %v", err))
        return
    }

    respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (a *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
    feedFollowIDStr := chi.URLParam(r, "feedFollowID")
    feedFollowID, err := uuid.Parse(feedFollowIDStr)

    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow ID: %v", err))
        return
    }

    errDel := a.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
        ID: feedFollowID,
        UserID: user.ID,
    })

    if errDel != nil {
        respondWithError(w, 500, fmt.Sprintf("Couldn't delete feed follow: %v", err))
        return
    }

    respondWithJSON(w, 204, struct{}{})
}
