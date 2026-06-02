package main

import (
	"sync"

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
	mu    sync.RWMutex
}
