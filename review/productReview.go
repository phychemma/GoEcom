package review

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
	ordermanagement "phyEcom.com/orderManagement"
	"phyEcom.com/utils"
)

type ReviewWithUsername struct {
	Comment  string
	Username string
	// Image    string
	ID     uint
	UserID uint
}

func review(db *gorm.DB, productID uint, limit, offset int) ([]ReviewWithUsername, error) {
	// var product models.Product
	var reviews []ReviewWithUsername
	err := db.Table("reviews").Joins("LEFT JOIN users ON users.id =reviews.user_id").Select("reviews.comment", "users.username", "reviews.id", "reviews.user_id").Order("reviews.created_at DESC").Limit(limit).Offset(offset).Where("product_id=?", productID).Scan(&reviews).Error
	if err != nil {
		log.Printf("Error fetching product reviews: %v", err)
		return reviews, err
	}
	return reviews, nil
}

// func reverseReview(messages []ReviewWithUsername) []ReviewWithUsername {
// 	reversed := make([]ReviewWithUsername, len(messages))
// 	for i, msg := range messages {
// 		reversed[len(messages)-1-i] = msg
// 	}
// 	return reversed
// }

func GetReviews(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productUint, parseErr := strconv.ParseUint(r.FormValue("productID"), 10, 34)
		if parseErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to parse string to uint")
			return
		}
		page, limit, PLerr := ordermanagement.ReturnPageLimit(r)
		if PLerr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "limit and page not set")
			return
		}
		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}

		// Calculate offset
		offset := (page - 1) * limit
		Reviews, err := review(db, uint(productUint), limit, offset)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to get reviews")
			return
		}
		json.NewEncoder(w).Encode(Reviews)
	}
}
