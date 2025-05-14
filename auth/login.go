package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"phyEcom.com/models"
	"phyEcom.com/utils"
)

var jwtSecret string

func init() {
	if os.Getenv("RENDER") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	jwtSecret = os.Getenv("JWT_SECRET")
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var dbUser models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		result := db.Model(&models.User{}).Where(&models.User{Email: user.Email}).First(&dbUser)
		if result.Error != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		if !dbUser.EmailVerified {
			utils.RespondWithError(w, http.StatusUnauthorized, "Email not verified")
			return
		}
		tokenString, tokenerr := getToken(dbUser.ID, dbUser.Username, dbUser.Email, dbUser.Role)
		if tokenerr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error generating token")
		}
		setTokenToCookie(w, "token", tokenString)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(models.LoggedIn{Loggedin: true}); err != nil {
			return
		}
	}
}
func Logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(time.Hour * 72 * -1),
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logged out"))
}
func AdminDashboard(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("Welcome to the admin dashboard"))
	if err != nil {
		return
	}
}

func GetUerSessionData(r *http.Request) (map[string]interface{}, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	tokenStr := cookie.Value
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}
	if !token.Valid {
		return nil, errors.New("expired token")
	}
	return claims, nil
}

func CheckIfAuthenticated(w http.ResponseWriter, r *http.Request) {
	data, err := GetUerSessionData(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error geting data from session")
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "encoding error")
		return
	}

}
