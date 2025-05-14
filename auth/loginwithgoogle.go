package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"phyEcom.com/models"
	"phyEcom.com/utils"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random"
)

func init() {
	if os.Getenv("RENDER") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		log.Fatal("Missing GOOGLE_CLIENT_ID environment variable")
	}

	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientSecret == "" {
		log.Fatal("Missing GOOGLE_CLIENT_SECRET environment variable")
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://127.0.0.1:333/callback",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}

func HandleCallbackFromGoogle(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Query().Get("state") != oauthStateString { // checking for  state by extracting
			utils.RespondWithError(w, http.StatusBadRequest, "State is invalid")
			return
		}
		code := r.URL.Query().Get("code")
		token, err := googleOauthConfig.Exchange(context.Background(), code) // extractimg the code variable for token use
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Code exchange failed")
			return
		}

		client := googleOauthConfig.Client(context.Background(), token) // getting the user details from token
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed getting user info")
			return
		}
		defer resp.Body.Close()

		var userInfo map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil { //marshal with map
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to parse user info")
			return
		}
		email := userInfo["email"].(string)
		name := userInfo["given_name"].(string)

		exist, returnedUser, err := emailExists2(db, email)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error checking email")
			return
		}
		if !returnedUser.EmailVerified { // delete if user exist and email is not verified
			db.Delete(&models.User{}, returnedUser.ID)
		}
		if !exist { // creating new user
			returnedUser.Role = "user"
			returnedUser.EmailVerified = true
			returnedUser.Email = email
			returnedUser.Username = name
			result := db.Create(&returnedUser)
			if result.Error != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error registering user")
				return
			}
		} else {
			if !returnedUser.EmailVerified {
				db.Model(&returnedUser).Updates(models.User{EmailVerified: true}) // since unverified email has been deleted this part is just meant to meant to update last login just a thought
			}
		}
		jwttoken, tokenErr := getToken(returnedUser.ID, name, email, "user") // set token
		if tokenErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}
		setTokenToCookie(w, "token", jwttoken)
		json.NewEncoder(w).Encode(userInfo)
	}
}
