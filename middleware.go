package main

import (
	"net/http"
    "fmt"

	"github.com/us0p/rss-aggregator/internal/auth"
	"github.com/us0p/rss-aggregator/internal/database"
)

type authHandler func (http.ResponseWriter, *http.Request, database.User)

func (a *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request){
        apiKey, err := auth.GetApiKey(r.Header)
    
        if err != nil {
            respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
            return
        }

        user, err := a.DB.GetUserByAPIKey(r.Context(), apiKey)
    
        if err != nil {
            respondWithError(w, 500, fmt.Sprintf("Couldn't get the user: %v", err))
            return
        }
        
        handler(w, r, user)
    }
}
