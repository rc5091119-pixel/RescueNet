package main

import (
	"context"
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
	ctx := context.TODO()
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

	cfg.hub.mu.Lock()
	if cfg.hub.Rooms[roomID] == nil {
		cfg.hub.Rooms[roomID] = make(map[*Client]bool)
	}

	cfg.hub.Rooms[roomID][client] = true
	cfg.hub.mu.Unlock()

	defer func() {
		cfg.hub.mu.Lock()
		delete(cfg.hub.Rooms[roomID], client)
		if len(cfg.hub.Rooms[roomID]) == 0 {
			delete(cfg.hub.Rooms, roomID)
		}

		log.Printf(
			"User %s disconnected from room %s",
			userID,
			roomID,
		)
		cfg.hub.mu.Unlock()
		conn.Close()
	}()
	type WSMessage struct {
		Type      string  `json:"type"`
		Content   string  `json:"content,omitempty"`
		Latitude  float64 `json:"latitude,omitempty"`
		Longitude float64 `json:"longitude,omitempty"`
	}

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var incoming WSMessage

		err = json.Unmarshal(data, &incoming)
		if err != nil {
			log.Println(err)
			continue
		}

		switch incoming.Type {

		case "chat":

			msg, err := cfg.db.CreateMessage(
				ctx,
				database.CreateMessageParams{
					RoomID:   roomID,
					SenderID: userID,
					Content:  incoming.Content,
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
			fullMsg, err := cfg.db.GetMessageByID(
				ctx,
				msg.ID,
			)

			if err != nil {
				log.Println(err)
				continue
			}

			payload, err := json.Marshal(fullMsg)
			if err != nil {
				log.Println(err)
				continue
			}

			cfg.hub.mu.RLock()

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
			cfg.hub.mu.RUnlock()

		case "location":

			payload, err := json.Marshal(incoming)
			if err != nil {
				continue
			}

			cfg.hub.mu.RLock()

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

			cfg.hub.mu.RUnlock()

		}
	}
}
