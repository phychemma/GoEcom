package profile

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gorm.io/gorm"
	"phyEcom.com/auth"
	"phyEcom.com/models"
	"phyEcom.com/utils"
)

func GetUserID(r *http.Request) (uint, error) {
	claims, err := auth.GetUerSessionData(r)
	if err != nil {
		return 0, err
	}
	floatUser, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("id does not exist")
	}
	return uint(floatUser), nil
}

func GetProfileDetailes(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		id, err := GetUserID(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error getting user id")
			return
		}
		db.Preload("Profile").First(&user, id)
		json.NewEncoder(w).Encode(user)

	}
}
func GetUserName(id uint, db *gorm.DB) string {
	var username string
	if err := db.Model(&models.User{}).Where(&models.User{ID: id}).Select("username").Scan(&username).Error; err != nil {
		log.Printf("error fetching username %v", err)
	}
	return username
}
