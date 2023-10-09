package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errResponse struct {
    Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
    if code > 499 {
        log.Println("Responding with code 5XX error:", msg)
    }

    respondWithJSON(w, code, errResponse{
        Error: msg,
    })
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    data, err := json.Marshal(payload) // Creates a JSON with the provided payload.
    if err != nil {
        log.Printf("Failed to marshal JSON response: %v, with error\n", payload, err)
        w.WriteHeader(500)
    }

    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(data)
}
