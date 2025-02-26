package review

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"phyEcom.com/middleware"
	"phyEcom.com/models"
	"phyEcom.com/profile"
	"phyEcom.com/utils"
)

type BroadcastReview struct {
	SenderID  uint   // ID of the sender
	Review    string // actual message
	ProductID uint   // product reference (if any)
}

type UsedBroadcastReview struct {
	models.Review
	Username string
}

type Client struct {
	UserID    uint // Unique identifier for the user
	Username  string
	ProductID uint
	Socket    *websocket.Conn
	Send      chan []byte
}
type ProductReviewManager struct {
	CLientInProductReview map[uint][]*Client
	Register              chan *Client
	Unregister            chan *Client
	mu                    sync.Mutex
	Broadcast             chan *UsedBroadcastReview
}

func NewProductReviewManager() *ProductReviewManager {
	return &ProductReviewManager{
		CLientInProductReview: make(map[uint][]*Client),
		Register:              make(chan *Client),
		Unregister:            make(chan *Client),
		Broadcast:             make(chan *UsedBroadcastReview),
	}
}

func (manager *ProductReviewManager) GetProductReviewConn(productID uint) *ProductReviewManager {
	_, ok := manager.CLientInProductReview[productID]
	if ok {
		return manager
	} else {
		manager.CLientInProductReview[productID] = []*Client{}
		return manager
	}
}

func (manager *ProductReviewManager) AddClient(c *Client) {
	clients, ok := manager.GetProductReviewConn(c.ProductID).CLientInProductReview[c.ProductID]
	if !ok {
		log.Println("product review does not exist addClient")
	}
	manager.CLientInProductReview[c.ProductID] = append(clients, c)
}

func (manager *ProductReviewManager) RemoveClient(c *Client) {
	clients, ok := manager.GetProductReviewConn(c.ProductID).CLientInProductReview[c.ProductID]
	if !ok {
		log.Println("product review does not exist reemoveClient")
	}
	for i, client := range clients {
		if c == client {
			manager.CLientInProductReview[c.ProductID] = append(clients[:i], clients[i+1:]...)
		}
	}
}

func (manager *ProductReviewManager) BroadCast(msg *UsedBroadcastReview) {
	var reviewWithUsername ReviewWithUsername
	toclients, ok := manager.GetProductReviewConn(msg.ProductID).CLientInProductReview[msg.ProductID]
	if !ok {
		log.Println("product review does not exist BroadCast")
	}
	for _, client := range toclients {
		reviewWithUsername.Comment = msg.Comment
		reviewWithUsername.Username = msg.Username
		reviewWithUsername.UserID = msg.UserID
		byt, jsonErr := json.Marshal(reviewWithUsername)
		if jsonErr != nil {
			log.Println("err converting to json")
		}
		client.Send <- byt
	}
}

func (manager *ProductReviewManager) Run() {
	for {
		select {
		case client := <-manager.Register:
			manager.mu.Lock()
			manager.AddClient(client)
			manager.mu.Unlock()
			log.Printf("User %d connected", client.UserID)

		case client := <-manager.Unregister:
			manager.mu.Lock()
			manager.RemoveClient(client)
			manager.mu.Unlock()

		case message := <-manager.Broadcast:
			manager.mu.Lock()
			manager.BroadCast(message)
			manager.mu.Unlock()
		}
	}
}

func (client *Client) WriteReview(manager *ProductReviewManager) {
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

func (client *Client) ReadReview(manager *ProductReviewManager, db *gorm.DB) {

	var user models.User
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
		var reviewModel models.Review
		var review BroadcastReview
		if err := json.Unmarshal(msg, &review); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}
		reviewModel.Comment = review.Review
		reviewModel.ProductID = review.ProductID
		reviewModel.UserID = review.SenderID
		if err := db.Create(&reviewModel).Error; err != nil {
			log.Println("error creating review")
			break
		}

		if err := db.Model(&models.User{}).Select("username").Where("id=?", reviewModel.UserID).Scan(&user).Error; err != nil {
			log.Printf("Error occured finding username %v", err)
			break
		}

		manager.Broadcast <- &UsedBroadcastReview{reviewModel, user.Username}

	}
}

func HandleReviewWebsocket(manager *ProductReviewManager, db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var productID uint64
	productid := vars["productID"]
	upgrader := websocket.Upgrader{
		CheckOrigin: middleware.CheckWebSocketOrigin,
	}
	if productid != "" {
		var err error
		productID, err = strconv.ParseUint(productid, 10, 32)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "id parsing error")
			return
		}
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
	log.Println(profile.GetUserName(userid, db))

	client := &Client{
		UserID:    userid,
		Username:  profile.GetUserName(userid, db),
		Socket:    conn,
		Send:      make(chan []byte),
		ProductID: uint(productID),
	}
	manager.Register <- client

	go client.ReadReview(manager, db)
	go client.WriteReview(manager)
}
