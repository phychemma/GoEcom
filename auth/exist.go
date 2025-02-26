package auth

import (
	"encoding/json"
	"net/http"

	"phyEcom.com/models"
	"phyEcom.com/utils"

	"gorm.io/gorm"
)

func CheckIfUserExist(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.UserExist
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error Decodig")
			return
		}
		exist, userE, err := userExists2(db, user.Username)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error checking for username")
			return
		}
		if !userE.EmailVerified && exist {
			user.Exist = false
		} else {
			user.Exist = exist
		}
		if err := json.NewEncoder(w).Encode(user); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error responding")
			return
		}
	}
}

func CheckIfEmailExist(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var email models.EmailExist
		if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error Decodig")
			return
		}
		exist, user, err := emailExists2(db, email.Email)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error checking for email")
			return
		}
		if !user.EmailVerified && exist {
			email.Exist = false
		} else {
			email.Exist = exist
		}
		if err := json.NewEncoder(w).Encode(email); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error responding")
			return
		}
	}
}
