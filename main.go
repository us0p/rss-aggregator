package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/us0p/rss-aggregator/internal/database"
)

type apiConfig struct {
    DB *database.Queries
}

func main() {
    godotenv.Load(".env")

    portString := os.Getenv("PORT")

    if portString == "" {
        log.Fatal("PORT is not defined in the environment")
    }

    dbURL := os.Getenv("DB_URL")

    if dbURL == "" {
        log.Fatal("DB_URL is not defined in the environment")
    }

    conn, err := sql.Open("postgres", dbURL)

    if err != nil {
        log.Fatal("Can't connect to the database", err)
    }

    db := database.New(conn)
    apiCfg := apiConfig{
        DB: db,
    }

    go startScraping(
        db,
        10,
        time.Minute,
    )

    router := chi.NewRouter()

    router.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"https://*", "http://*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"*"},
        ExposedHeaders: []string{"Link"},
        AllowCredentials: false,
        MaxAge: 300,
    }))

    v1Router := chi.NewRouter()
    v1Router.Get("/healthz", handlerReadiness)
    v1Router.Post("/users", apiCfg.handlerCreateUser)
    v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
    v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
    v1Router.Get("/feeds", apiCfg.handleGedFeeds)
    v1Router.Post("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
    v1Router.Get("/feed_follow", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
    v1Router.Delete("/feed_follow/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))
    v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

    router.Mount("/v1", v1Router)

    server := &http.Server{
        Handler: router,
        Addr: fmt.Sprintf(":%s", portString),
    }

    servErr := server.ListenAndServe()

    if servErr != nil {
        log.Fatal(servErr)
    }
}
