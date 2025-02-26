package websocket

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
	"phyEcom.com/auth"
	"phyEcom.com/models"
	"phyEcom.com/profile"
	"phyEcom.com/utils"
)

func FetcUsers(db *gorm.DB) http.HandlerFunc {
	type sending struct {
		Users  []models.User
		UserID uint
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userid, uerr := profile.GetUserID(r)
		if uerr != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Error getting user id")
			return
		}
		users, err := auth.FetcUsers(userid, db)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Error getting other users")
			return
		}
		json.NewEncoder(w).Encode(sending{Users: users, UserID: userid})
	}
}
