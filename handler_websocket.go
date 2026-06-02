package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rc5091119-pixel/rescuenet/internal/auth"
	"github.com/rc5091119-pixel/rescuenet/internal/database"
)

func (cfg *apiConfig) handlerWebsocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	log.Println("Token:", token)

	userID, err := auth.ValidateJWT(
		token,
		cfg.jwtSecret,
	)
	log.Println("UserID:", userID)

	if err != nil {
		http.Error(
			w,
			"Unauthorized",
			http.StatusUnauthorized,
		)
		return
	}
	roomIDstr := r.PathValue("roomID")
	roomID, err := uuid.Parse(roomIDstr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid RoomId", err)
		return
	}
	err = cfg.VerifyRoomMember(r.Context(), roomID, userID)
	if err != nil {
		respondWithError(
			w,
			http.StatusForbidden,
			"You are not a member of this room",
			err,
		)
		return
	}
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	client := &Client{
		Conn:   conn,
		UserID: userID,
		RoomID: roomID,
	}
	if cfg.hub.Rooms[roomID] == nil {
		cfg.hub.Rooms[roomID] = make(map[*Client]bool)
	}
	cfg.hub.Rooms[roomID][client] = true
	log.Println("WebSocket connected")

	defer func() {
		delete(cfg.hub.Rooms[roomID], client)

		log.Printf(
			"User %s disconnected from room %s",
			userID,
			roomID,
		)

		conn.Close()
	}()
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		content := string(data)

		log.Printf(
			"Saving message. Room=%s User=%s Content=%s",
			roomID,
			userID,
			content,
		)

		msg, err := cfg.db.CreateMessage(
			r.Context(),
			database.CreateMessageParams{
				RoomID:   roomID,
				SenderID: userID,
				Content:  content,
			},
		)

		if err != nil {
			log.Printf("DB Error : %v", err)
			continue
		}
		log.Printf(
			"Saved message %s",
			msg.ID,
		)

		payload, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf(
			"Broadcasting message %s to %d clients",
			msg.ID,
			len(cfg.hub.Rooms[roomID]),
		)

		for c := range cfg.hub.Rooms[roomID] {
			if c == client {
				continue
			}

			err := c.Conn.WriteMessage(
				websocket.TextMessage,
				payload,
			)

			if err != nil {
				log.Println(err)
			}
		}
	}
}
