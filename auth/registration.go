package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"phyEcom.com/models"
	"phyEcom.com/utils"
)

func userIsAdmin(email string) bool {
	boolReturn := false
	admin := os.Getenv("ADMIN")
	arr := strings.Split(admin, ",")
	for _, v := range arr {
		if v == email {
			boolReturn = true
			break
		}
	}
	return boolReturn
}
func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid user request payload")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
			return
		}
		verificationCode := getCode()
		user.Password = string(hashedPassword)

		if userIsAdmin(user.Email) {
			user.Role = "admin"
		} else {
			user.Role = "user"
		}

		user.VerificationCode = verificationCode
		if err := sendAuthMail(verificationCode, user.Email); err != nil {
			//utils.RespondWithError(w, http.StatusInternalServerError, "Error sending Email")
			//return
		}
		log.Println((verificationCode))
		emailExist, userEmailExist, err := emailExists2(db, user.Email)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error checking for email")
			return
		}
		if !userEmailExist.EmailVerified && emailExist {
			result := db.Delete(&models.User{}, userEmailExist.ID)
			if result.Error != nil {
				log.Println()
			}
		}
		user.Username = strings.Split(user.Email, "@")[0]

		userExist, _, userErr := userExists2(db, user.Username)
		if userErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error checking for username")
			return
		}
		userid := 0
		for userExist {
			userid += 1
			user.Username = fmt.Sprintf(`%s%d`, user.Username, userid)
			userExist, _, userErr = userExists2(db, user.Username)
			if userErr != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error checking for username")
				return
			}
		}

		result := db.Create(&user)
		if result.Error != nil {
			log.Print(result.Error)
			utils.RespondWithError(w, http.StatusInternalServerError, "Error registering user")
			return
		}

		user.Password = ""
		user.Role = ""

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(user); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error responding")
			return
		}
	}
}

func VerifyCode(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var code models.VerifyCode
		var user models.User
		var count int64
		if err := json.NewDecoder(r.Body).Decode(&code); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Could not decode data")
			return
		}
		log.Printf("email %s code %s", code.Email, code.Code)

		result := db.Model(&models.User{}).Where(&models.User{Email: code.Email, VerificationCode: code.Code}).Count(&count)
		log.Println(result.Debug())
		if result.Error != nil {
			log.Println(result.Error)
			return
		}

		if count == 0 {
			utils.RespondWithError(w, http.StatusInternalServerError, "No Email with such verification code")
			return
		}
		db.Model(&models.User{}).Where(&models.User{Email: code.Email, VerificationCode: code.Code}).First(&user).Updates(&models.User{EmailVerified: true, VerificationCode: "used"})
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(code); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error responding")
			return
		}
	}
}

func getCode() string {
	return fmt.Sprint(time.Now().Nanosecond())[3:]
}

func sendAuthMail(code, email string) error {
	html := fmt.Sprintf(`<h2>Email authentication from <strong>Vidstream </strong></h><br/> <b>Your reg code is <i> %s</i><b/>
	<strong> Ignore if you did not make this request </strong>
	`, code)
	if err := Mailing("Email Authentication VidStream", &html, email); err != nil {
		return err
	}
	return nil
}
