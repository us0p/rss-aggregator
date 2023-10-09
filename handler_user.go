package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/us0p/rss-aggregator/internal/database"
)

type parameters struct {
    Name string `json:"name"`
}

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)

    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
        return
    }

    user, err := a.DB.CreateUser(r.Context(), database.CreateUserParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: params.Name,
    })

    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
        return
    }

    respondWithJSON(w, 201, databaseUserToUser(user))
}

func (a *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
    respondWithJSON(w, 200, databaseUserToUser(user))
}

func (a *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
    posts, err := a.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
        UserID: user.ID,
        Limit: 10,
    })

    if err != nil {
        respondWithError(w, 500, fmt.Sprintf("Couldn't get posts for user: %v", err))
        return
    }

    respondWithJSON(w, 200, databasePostsToPosts(posts))
}
