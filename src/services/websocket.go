package services

import (
	"ankasa-be/src/models"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	Rooms map[int]*Room
	Mutex sync.Mutex
}

type Room struct {
	Clients    map[*websocket.Conn]bool
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Broadcast  chan models.Message
}

var hub *Hub

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[int]*Room),
	}
}

func InitHub() {
	hub = NewHub()
}

func GetHub() *Hub {
	return hub
}

func (r *Room) Run() {
	for {
		select {
		case conn := <-r.Register:
			r.Clients[conn] = true
		case conn := <-r.Unregister:
			delete(r.Clients, conn)
		case msg := <-r.Broadcast:
			for conn := range r.Clients {
				_ = conn.WriteJSON(msg)
			}
		}
	}
}

func (h *Hub) GetRoom(roomID int) *Room {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	room, exists := h.Rooms[roomID]
	if !exists {
		room = &Room{
			Clients:    make(map[*websocket.Conn]bool),
			Register:   make(chan *websocket.Conn),
			Unregister: make(chan *websocket.Conn),
			Broadcast:  make(chan models.Message),
		}
		h.Rooms[roomID] = room
		go room.Run()
	}
	return room
}
