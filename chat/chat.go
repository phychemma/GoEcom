package chat

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
	"phyEcom.com/models"
	ordermanagement "phyEcom.com/orderManagement"
	"phyEcom.com/product"
	"phyEcom.com/profile"
	"phyEcom.com/utils"
)

type ProductWithLastMessage struct {
	ProductID   uint
	ProductName string
	LastMessage string
	MessageTime time.Time
}

func GetTotalProductCount(db *gorm.DB) (int64, error) {
	var count int64
	if err := db.Model(&models.Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func FetchProductsWithChats(db *gorm.DB, userType string, page, limit int) ([]ProductWithLastMessage, error) {
	var results []ProductWithLastMessage
	var query string
	offset := (page - 1) * limit // Calculate the offset based on the page

	// Raw SQL query with pagination
	query1 := `
        SELECT p.id AS product_id, p.name AS product_name, c.message AS last_message, c.created_at AS message_time
        FROM products p
        LEFT JOIN chats c ON p.id = c.product_id
		`
	query2 := `
		LEFT JOIN users u ON p.user_id = u.id
		`
	query3 := `
        WHERE c.id = (
            SELECT id FROM chats
            WHERE product_id = p.id
            ORDER BY created_at DESC
            LIMIT 1
        )
        ORDER BY c.created_at DESC
        LIMIT ? OFFSET ?
    `
	if userType == "admin" {
		query = query1 + query3

	} else if userType == "user" {
		query = query1 + query2 + query3
	} else {
		return results, errors.New("user type unknowm ")
	}
	// Execute the query
	if err := db.Raw(query, limit, offset).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func GetProducts(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		userID, uerr := profile.GetUserID(r)
		if uerr != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "user not authorized")
			return
		}
		if roleErr := db.Model(&models.User{}).Find(&user, userID).Error; roleErr != nil {
			log.Println(roleErr)
			utils.RespondWithError(w, http.StatusUnauthorized, "unable to fetch user role ")
			return
		}
		page, limit, PLerr := ordermanagement.ReturnPageLimit(r)
		if PLerr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "limit and page not set")
			return
		}
		product, productErr := FetchProductsWithChats(db, user.Role, page, limit)
		if productErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch data")
			return
		}
		utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"product": product})

	}
}
func ChattedProduct(db *gorm.DB, userID uint) ([]product.ProductOneChat, error) {
	var product []product.ProductOneChat
	query := `
        SELECT 
            p.id AS product_id,
            p.name AS product_name,
            p.brand AS product_brand,
			pic.url as product_pic,
            COALESCE(c.message, '') AS last_message, -- Default to empty string if no chat
            COALESCE(c.created_at, '0001-01-01 00:00:00') AS message_time -- Default to zero time if no chat
        FROM 
			products p
        LEFT JOIN 
			chats c
            
        ON 
            p.id = c.product_id 
            AND c.id = (
                SELECT id 
                FROM chats 
                WHERE product_id = p.id 
                ORDER BY created_at DESC 
                LIMIT 1
            )
		LEFT JOIN 
			pictures pic
		ON 
				p.id = pic.product_id
			AND pic.id = (
				SELECT id FROM pictures
				WHERE product_id = p.id
				ORDER BY created_at DESC
				LIMIT 1
			)
        WHERE 
            c.user_id = ?
    `
	if err := db.Raw(query, userID).Scan(&product).Error; err != nil {
		return product, errors.New("list of chatted product not found ")
	}
	return product, nil
}

func FetchUserChattedProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := profile.GetUserID(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "user not autenticated")
			return
		}
		products, pErr := ChattedProduct(db, userID)
		if pErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "product not found")
			return
		}
		json.NewEncoder(w).Encode(products)
	}
}

type ChatForFrontEnd struct {
	ChatID    uint      `json:"chat_id"`
	Message   string    `json:"message"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	Read      bool      `json:"read"`
}

func GetUserProductChat(db *gorm.DB, productID, UserID uint, page int) ([]ChatForFrontEnd, error) {
	var results []ChatForFrontEnd
	offset := (page - 1) * 15
	if err := db.Debug().Table("chats").
		Select("id AS chat_id, message, is_admin, `read` , created_at").
		Where("user_id= ? && product_id=?", UserID, productID).
		Order("id DESC").
		Limit(15).
		Offset(offset).
		Scan(&results).Error; err != nil {
		return results, err
	}
	return results, nil
}
func GetChats(db *gorm.DB) http.HandlerFunc {
	type RequestData struct {
		Page      int  `json:"chatPage"`
		UserID    uint `json:"userID"`
		ProductID uint `json:"productID"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData RequestData
		json.NewDecoder(r.Body).Decode(&requestData)
		results, err := GetUserProductChat(db, requestData.ProductID, requestData.UserID, requestData.Page)
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusBadRequest, "err getting user chat")
			return
		}
		json.NewEncoder(w).Encode(results)
	}
}

// func GetChats(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userID, uerr := profile.GetUserID(r)
// 		page, limit, PLerr := ordermanagement.ReturnPageLimit(r)
// 		if PLerr != nil {
// 			utils.RespondWithError(w, http.StatusInternalServerError, "limit and page not set")
// 			return
// 		}
// 		if page <= 0 {
// 			page = 1
// 		}
// 		if limit <= 0 {
// 			limit = 10
// 		}
// 		offset := (page - 1) * limit

// 	}
// }
