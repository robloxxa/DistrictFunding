package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/robloxxa/DistrictFunding/internal/auth"
	"github.com/robloxxa/DistrictFunding/pkg/jwtauth"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	log.Println("Connecting to database")
	pool, err := pgxpool.New(context.Background(),
		fmt.Sprintf("postgres://postgres:%s@%s/auth_db",
			os.Getenv("AUTH_POSTGRES_PASSWORD"),
			os.Getenv("AUTH_POSTGRES_HOST"),
		))
	if err != nil {
		fmt.Printf("Unable to create connection pool: %v", err)
		return
	}

	defer pool.Close()

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Printf("Unable to acquire database connection")
		return
	}
	conn.Release()
	// Initialize JWT auth

	log.Println("Connected to Database")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	jwt_secret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatalln("No jwt secret variable")
	}
	ja := jwtauth.New("HS256", []byte(jwt_secret))

	r.Mount("/", auth.NewController(pool, ja))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
