package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

type apiConfig struct {
	db        *database.Queries
	jwtSecret string
}

func main() {
	const port = "8080"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("You mush set DBURL")
	}
	dbconn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbconn)

	secretKey := os.Getenv("JWT_SECRET")
	mux := http.NewServeMux()
	apiconfig := apiConfig{
		db:        dbQueries,
		jwtSecret: secretKey,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("RescueNet API running 🚀"))
	})
	mux.HandleFunc("/api/users", apiconfig.handlerCreateUsers)
	mux.HandleFunc("/api/login", apiconfig.handlerLoginUsers)
	mux.Handle("/api/test",
    apiconfig.AuthMiddleware(http.HandlerFunc(apiconfig.handlerTestProtected)))
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
