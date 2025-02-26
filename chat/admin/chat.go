package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"phyEcom.com/chat"
	"phyEcom.com/profile"
	"phyEcom.com/utils"
)

func FectchAllChattedUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var users []struct {
			UserID    uint      `json:"id"`
			Username  string    `json:"username"`
			Email     string    `json:"email"`
			LastChat  time.Time `json:"last_chat"`  // Timestamp of the latest chat
			ChatCount int       `json:"chat_count"` // Number of chats by the user
		}
		userid, uerr := profile.GetUserID(r)
		if uerr != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "userID not found")
			return
		}
		err := db.Table("users").
			Select(`users.id AS user_id, users.username, users.email,
					MAX(chats.created_at) AS last_chat, COUNT(chats.id) AS chat_count`).
			Joins("JOIN chats ON chats.user_id = users.id").
			Group("users.id, users.username, users.email").
			Where("users.id != ?", userid).
			Order("last_chat DESC"). // Order by the latest chat timestamp
			Scan(&users).Error

		if err != nil {
			utils.RespondWithError(w, http.StatusNotFound, "Error fetching data")
			return
		}

		json.NewEncoder(w).Encode(users)

	}
}

func GetProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID uint64
		vars := mux.Vars(r)
		userStrVal := vars["userID"]
		if userStrVal != "" {
			var err error
			userID, err = strconv.ParseUint(userStrVal, 10, 32)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "error coverting productID to an integer")
				return
			}
		}
		products, pErr := chat.ChattedProduct(db, uint(userID))
		if pErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "product not found")
			return
		}
		json.NewEncoder(w).Encode(products)
	}
}
