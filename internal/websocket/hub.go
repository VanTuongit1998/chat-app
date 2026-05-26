package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"chat-app/internal/model"
	"chat-app/internal/service"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients     map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	broadcast   chan []byte
	redis       *redis.Client
	pubsub      *redis.PubSub
	chatUsecase *service.ChatUsecase
	mu          sync.RWMutex
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

func NewHub(chatUsecase *service.ChatUsecase, redisClient *redis.Client) *Hub {
	pubsub := redisClient.Subscribe(context.Background(), "chat_messages")
	return &Hub{
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		redis:       redisClient,
		pubsub:      pubsub,
		chatUsecase: chatUsecase,
	}
}

func (h *Hub) Run(ctx context.Context) {
	go h.listenRedis(ctx)

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			h.sendMessage(message)
			h.mu.Unlock()
		case <-ctx.Done():
			h.pubsub.Close()
			return
		}
	}
}

func (h *Hub) listenRedis(ctx context.Context) {
	channel := h.pubsub.Channel()
	for {
		select {
		case msg, ok := <-channel:
			if !ok {
				return
			}
			h.broadcast <- []byte(msg.Payload)
		case <-ctx.Done():
			return
		}
	}
}

func (h *Hub) publishMessage(msg *model.Message) error {
	return h.chatUsecase.PublishMessage(context.Background(), msg)
}

func (h *Hub) sendMessage(message []byte) {
	var msg model.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("[Hub] Unmarshal failed: %v, broadcasting raw message", err)
		for client := range h.clients {
			h.sendToClient(client, message)
		}
		return
	}

	log.Printf("[Hub] Broadcasting message from %s to %s", msg.Sender, msg.To)
	count := 0
	for client := range h.clients {
		if msg.To != "" && client.username != msg.Sender && client.username != msg.To {
			continue
		}
		h.sendToClient(client, message)
		count++
	}
	log.Printf("[Hub] Sent to %d clients", count)
}

func (h *Hub) sendToClient(client *Client, message []byte) {
	select {
	case client.send <- message:
	default:
		close(client.send)
		delete(h.clients, client)
	}
}
