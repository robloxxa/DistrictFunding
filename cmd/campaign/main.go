package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/robloxxa/DistrictFunding/internal/campaign"
	"github.com/robloxxa/DistrictFunding/pkg/jwtauth"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	log.Println("Connecting to database")
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s/test", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST")))
	if err != nil {
		fmt.Printf("Unable to create connection pool: %v", err)
		return
	}

	defer pool.Close()

	_, err = pool.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Unable to acquire database connection")
		return
	}

	// Initialize JWT auth

	log.Println("Connected to Database")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	ja := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")))

	r.Mount(`/`, campaign.NewController(pool, ja))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
