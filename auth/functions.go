package auth

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"phyEcom.com/models"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

func getToken(id uint, username, email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"email":    email,
		"role":     role,
		"id":       id,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func setTokenToCookie(w http.ResponseWriter, name, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}

// func userExists(db *gorm.DB, username string) (bool, error) {
// 	var count int64
// 	result := db.Where("username = ?", username).Count(&count)
// 	if result.Error != nil {
// 		if result.Error == gorm.ErrRecordNotFound {
// 			return false, nil
// 		}
// 		return false, result.Error
// 	}
// 	return true, nil
// }

func userExists2(db *gorm.DB, username string) (bool, models.User, error) {
	var user models.User
	result := db.Select([]string{"id", "email_verified"}).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, user, nil
		}
		return false, user, result.Error
	}
	return true, user, nil
}

//	func emailExists(db *gorm.DB, email string) (bool, error) {
//		var user models.User
//		result := db.Where("email = ?", email).First(&user)
//		if result.Error != nil {
//			if result.Error == gorm.ErrRecordNotFound {
//				return false, nil
//			}
//			return false, result.Error
//		}
//		return true, nil
//	}
func emailExists2(db *gorm.DB, email string) (bool, models.User, error) {
	var user models.User
	result := db.Select([]string{"id", "email_verified"}).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, user, nil
		}
		return false, user, result.Error
	}
	return true, user, nil
}

func GenerateSKU(name string) string {
	// Create a simple SKU using the product name and current timestamp
	// You can customize this logic as needed
	namePart := strings.ToUpper(strings.ReplaceAll(name, " ", ""))
	timestampPart := time.Now().UnixNano()
	return fmt.Sprintf("%s-%d", namePart, timestampPart)
}

func NormalizePath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

func FetcUsers(userid uint, db *gorm.DB) ([]models.User, error) {
	var users []models.User
	if err := db.Where("id != ?", userid).Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}
