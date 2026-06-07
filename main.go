package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

type apiConfig struct {
	db        *database.Queries
	jwtSecret string
	hub       *Hub
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

	hub := &Hub{
		Rooms: make(map[uuid.UUID]map[*Client]bool),
	}

	mux := http.NewServeMux()
	apiconfig := apiConfig{
		db:        dbQueries,
		jwtSecret: secretKey,
		hub:       hub,
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("RescueNet API running 🚀"))
	})

	mux.HandleFunc("POST /api/users", apiconfig.handlerCreateUsers)

	mux.HandleFunc("POST /api/login", apiconfig.handlerLoginUsers)

	mux.Handle(
		"GET /api/test",
		apiconfig.AuthMiddleware(http.HandlerFunc(apiconfig.handlerTestProtected)),
	)

	mux.Handle(
		"POST /api/location",
		apiconfig.AuthMiddleware(http.HandlerFunc(apiconfig.handlerUpdateLocation)),
	)

	mux.Handle(
		"POST /api/alerts",
		apiconfig.AuthMiddleware(http.HandlerFunc(apiconfig.handlerCreateAlerts)),
	)

	mux.Handle(
		"GET /api/notifications",
		apiconfig.AuthMiddleware(
			http.HandlerFunc(apiconfig.handlerGetNotifications),
		),
	)

	mux.Handle(
		"POST /api/alerts/{id}/accept",
		apiconfig.AuthMiddleware(http.HandlerFunc(apiconfig.handlerAcceptAlert)),
	)
	mux.Handle(
		"GET /api/my-rooms",
		apiconfig.AuthMiddleware(
			http.HandlerFunc(apiconfig.handlerGetUserRooms),
		),
	)

	mux.Handle(
		"GET /api/rooms/{roomID}/info",
		apiconfig.AuthMiddleware(http.HandlerFunc(apiconfig.handlerGetRoomInfo)),
	)

	mux.Handle(
		"POST /api/rooms/{roomID}/messages",
		apiconfig.AuthMiddleware(http.HandlerFunc(apiconfig.handlerCreateMessage)),
	)

	mux.Handle(
		"GET /api/rooms/{roomID}/messages",
		apiconfig.AuthMiddleware(http.HandlerFunc(apiconfig.handlerGetMessage)),
	)
	mux.Handle(
		"GET /api/rooms/{roomID}/locations",
		apiconfig.AuthMiddleware(http.HandlerFunc(apiconfig.handlerGetRoomMemberLocations)),
	)

	mux.HandleFunc(
		"GET /ws/rooms/{roomID}",
		apiconfig.handlerWebsocket,
	)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMiddleware(mux),
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
