package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"phyEcom.com/chat"
	"phyEcom.com/middleware"
	"phyEcom.com/models"
	"phyEcom.com/profile"
	"phyEcom.com/utils"
)

var ADMIN uint = 3

type Client struct {
	ID     uint
	UserID uint // Unique identifier for the user
	Socket *websocket.Conn
	Send   chan []byte
}

type WebSocketManager struct {
	Clients    map[uint]*Client // Map of UserID to Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan toFrontEnd // Modified to include targeted messages
	mu         sync.Mutex
}

type toFrontEnd struct {
	chat.ChatForFrontEnd
	TO        uint
	From      uint
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
}

type BroadcastMessage struct {
	Type      string // message, typing, etc.
	UserID    uint   // ID of the sender
	ProductID uint   // product reference (if any)
	Role      string
	AdminID   uint
	ChatID    uint      `json:"chat_id"`
	Message   string    `json:"message"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	Read      bool      `json:"read"`
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan toFrontEnd),
	}
}

func (manager *WebSocketManager) Run() {
	for {
		select {
		case client := <-manager.Register:
			manager.mu.Lock()
			manager.Clients[client.UserID] = client
			manager.mu.Unlock()
			log.Printf("User %d connected", client.UserID)

		case client := <-manager.Unregister:
			manager.mu.Lock()
			if _, ok := manager.Clients[client.UserID]; ok {
				close(client.Send)
				delete(manager.Clients, client.UserID)
				log.Printf("User %d disconnected", client.UserID)
			}
			manager.mu.Unlock()

		case message := <-manager.Broadcast:
			manager.mu.Lock()
			byt, jerr := json.Marshal(message)
			if jerr != nil {
				log.Print("couldn't marshal")
			}
			if client, ok := manager.Clients[message.From]; ok { // this sends to you
				client.Send <- byt
			} else {
				log.Printf("User %d is not connected", message.From)
			}
			if client, ok := manager.Clients[message.TO]; ok { // this sends to your receiver
				client.Send <- byt
			} else {
				log.Printf("User %d is not connected", message.TO)
			}

			manager.mu.Unlock()
		}
	}
}

func (client *Client) readMessages(manager *WebSocketManager, db *gorm.DB, userID uint) {
	var toFrontEnd toFrontEnd
	defer func() {
		manager.Unregister <- client
		client.Socket.Close()
	}()

	for {
		_, msg, err := client.Socket.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		var broadcastMsg BroadcastMessage

		if err := json.Unmarshal(msg, &broadcastMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			break
		}

		// broadcastMsg.UserID = client.UserID // Ensure the sender ID is accurate
		databaseMsg := models.Chat{Message: broadcastMsg.Message, UserID: broadcastMsg.UserID, ProductID: broadcastMsg.ProductID}
		if broadcastMsg.Role == "admin" {
			databaseMsg.IsAdmin = true
			databaseMsg.AdminID = broadcastMsg.AdminID
			toFrontEnd.TO = broadcastMsg.UserID
		} else if broadcastMsg.Role == "user" {
			toFrontEnd.TO = ADMIN
		} else {
			log.Println("user type not found")
		}
		if err := db.Model(&models.Chat{}).Create(&databaseMsg).Error; err != nil {
			log.Printf("Error creating message: %v", err)
			break
		}
		toFrontEnd.ChatID = databaseMsg.ID
		toFrontEnd.Message = databaseMsg.Message
		toFrontEnd.IsAdmin = databaseMsg.IsAdmin
		toFrontEnd.CreatedAt = databaseMsg.CreatedAt
		toFrontEnd.Read = databaseMsg.Read
		toFrontEnd.ProductID = databaseMsg.ProductID
		toFrontEnd.UserID = databaseMsg.UserID
		toFrontEnd.From = userID
		manager.Broadcast <- toFrontEnd
	}
}

func (client *Client) writeMessage(manager *WebSocketManager) {
	defer func() {
		manager.Unregister <- client
		client.Socket.Close()
	}()

	for message := range client.Send {
		err := client.Socket.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}

// HandleWebSocket upgrades HTTP to WebSocket and registers the client
func HandleWebSocket(manager *WebSocketManager, db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: middleware.CheckWebSocketOrigin,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	userid, uerr := profile.GetUserID(r)
	if uerr != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "userid not found")
		return
	}

	client := &Client{
		UserID: userid,
		Socket: conn,
		Send:   make(chan []byte),
	}

	manager.Register <- client

	go client.readMessages(manager, db, userid)
	go client.writeMessage(manager)
}
