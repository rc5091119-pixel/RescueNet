package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID uuid.UUID
	RoomID uuid.UUID
}

type Hub struct {
	Rooms map[uuid.UUID]map[*Client]bool
}
