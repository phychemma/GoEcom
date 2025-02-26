package ordermanagement

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"phyEcom.com/models"
	"phyEcom.com/profile"
	"phyEcom.com/utils"
)

// getData(`${backEndPath}/orders/made?page=1&limit=10`),
// getData(`${backEndPath}/orders/received?page=1&limit=10`),
// getData(`${backEndPath}/orders/to-confirm`),

func FetchOrdersMade(db *gorm.DB, userID uint, page, limit int) ([]models.Order, bool) {
	var orders []models.Order
	offset := (page - 1) * limit
	result := db.Where("user_id = ?", userID).Order("created_at DESC").Preload("Product").Offset(offset).Limit(limit).Find(&orders)

	hasMore := result.RowsAffected > int64(limit)
	return orders, hasMore
}

func ReturnPageLimit(r *http.Request) (int, int, error) {
	var page, limit int
	vars := mux.Vars(r)
	pageStr := vars["page"]
	limitStr := vars["limit"]
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return 0, 0, errors.New("invalid page parameter")
		}
	}
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, errors.New("invalid limit parameter")
		}
	}
	return page, limit, nil
}

func GetOrderMade(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, limit, err := ReturnPageLimit(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "parameters not found")
			return
		}
		id, err := profile.GetUserID(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "user not found")
			return
		}
		orders, hasmore := FetchOrdersMade(db, id, page, limit)
		utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
			"order":   orders,
			"hasmore": hasmore,
		})
	}
}

func FetchOrdersReceived(db *gorm.DB, userID uint, page, limit int) ([]models.Order, bool) {
	var orders []models.Order
	offset := (page - 1) * limit
	// &models.Order{Product: models.Product{UserID: userID}, Seller: false}
	result := db.Joins("JOIN products ON products.id =orders.product_id AND products.user_id =? AND seller = ?", userID, 0).Order("created_at DESC").Preload("Product").Offset(offset).Limit(limit).Find(&orders)
	hasMore := result.RowsAffected > int64(limit)
	return orders, hasMore
}

func GetOrderReceived(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := profile.GetUserID(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "user not found")
			return
		}
		page, limit, err := ReturnPageLimit(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "parameters not found")
			return
		}
		arrOfOrder, hasMore := FetchOrdersReceived(db, userID, page, limit)
		utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"order": arrOfOrder, "hasMore": hasMore})
	}
}

func FetchOrdersToBeConfiredByBuyer(db *gorm.DB, userID uint, page, limit int) ([]models.Order, bool) {
	var orders []models.Order
	offset := (page - 1) * limit
	result := db.Model(&models.Order{}).Where(&models.Order{UserID: userID, Seller: true}).Where("buyer =?", 0).Order("created_at DESC").Preload("Product").Offset(offset).Limit(limit).Find(&orders)
	hasMore := result.RowsAffected > int64(limit)
	return orders, hasMore
}

func GetConfirmReceivedOrder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := profile.GetUserID(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "user not found")
			return
		}
		page, limit, err := ReturnPageLimit(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "parameters not found")
			return
		}
		arrOfOrder, hasMore := FetchOrdersToBeConfiredByBuyer(db, userID, page, limit)
		utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"order": arrOfOrder, "hasMore": hasMore})
	}
}

func ConfirmDeliveredOrder(db *gorm.DB) http.HandlerFunc { // confirms from seller
	type front struct {
		Delivered bool
		OrderID   uint
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var f front
		var status string
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "bad data keys")
			return
		}
		if f.Delivered {
			status = "Confirmation"
		} else {
			status = "pending"
		}
		db.Where(&models.Order{ID: f.OrderID}).Updates(&models.Order{Seller: f.Delivered, Status: status})
		json.NewEncoder(w).Encode(f)
	}
}

func ConfirmDelivered(db *gorm.DB) http.HandlerFunc { //confirmation from buyers
	type front struct {
		Delivered bool
		OrderID   uint
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var f front
		var status string
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "bad data keys")
			return
		}
		if f.Delivered {
			status = "Confirm"
		} else {
			status = "Confirmation"
		}
		db.Where(&models.Order{ID: f.OrderID}).Updates(&models.Order{Buyer: f.Delivered, Status: status})
		json.NewEncoder(w).Encode(f)
	}
}
