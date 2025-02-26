package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"phyEcom.com/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var Front string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Front = os.Getenv("FRONTEND")
}

var jwtSecret string

func EnableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", Front)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		next.ServeHTTP(w, r)

	})

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "username", (*claims)["username"])
		ctx = context.WithValue(ctx, "role", (*claims)["role"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(requiredRole string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value("role").(string)
			if role != requiredRole {
				utils.RespondWithError(w, http.StatusForbidden, "Forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CheckWebSocketOrigin(r *http.Request) bool { // phychemma this fuction is used for checking websocket origin
	// Grab the request origin
	origin := r.Header.Get("Origin")
	switch origin {
	case Front:
		return true
	default:
		return false

	}
	//return true
}
