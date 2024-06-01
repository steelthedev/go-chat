package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type WebSocketServer struct {
	id        string
	clients   map[*websocket.Conn]bool
	broadcast chan *Message
}

func NewWebSocket() *WebSocketServer {
	return &WebSocketServer{
		id:        uuid.New().String(),
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan *Message),
	}
}

func (s *WebSocketServer) HandleConnections(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (s *WebSocketServer) HandleWebSocket(ctx *websocket.Conn) {

	// Register a new Client
	s.clients[ctx] = true
	defer func() {
		delete(s.clients, ctx)
		ctx.Close()
	}()

	for {
		_, msg, err := ctx.ReadMessage()
		if err != nil {
			log.Println("Read Error:", err)
			break
		}

		// send the message to the broadcast channel
		log.Println(string(msg))
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Fatalf("Error Unmarshalling")
		}
		message.ClientName = s.id

		s.broadcast <- &message
	}
}

func (s *WebSocketServer) HandleMessages() {
	for {
		msg := <-s.broadcast

		// Send the message to all Clients

		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, getMessageTemplate(msg))
			if err != nil {
				log.Printf("Write  Error: %v ", err)
				client.Close()
				delete(s.clients, client)
			}

		}

	}
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("views/message.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	// Render the template with the message as data.
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, msg)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	return renderedMessage.Bytes()
}
