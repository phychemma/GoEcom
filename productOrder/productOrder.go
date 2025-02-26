package productorder

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
	"phyEcom.com/models"
	"phyEcom.com/profile"
	"phyEcom.com/utils"
)

// fd.append("userID", userID)
// fd.append("productID", product.ID)
// fd.append("quantity", quantity)
// fd.append("price", price)
// func fetchOrder(db *gorm.DB, productID uint) ([]models.Order, error) {
// 	var orders []models.Order
// 	if err := db.Find(&orders, productID).Error; err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }

func fetchUserOrder(db *gorm.DB, productID, userID uint) (models.Order, error) {
	var order models.Order
	if err := db.Where(&models.Order{UserID: userID, ProductID: productID}).First(&order).Error; err != nil {
		return order, err
	}
	return order, nil
}

func FetchProductOrder(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func FetchUserProductOrder(db *gorm.DB) http.HandlerFunc {
	type front struct {
		ProductID uint
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var f front
		if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Could not Decode")
			return
		}
		fmt.Println(f.ProductID)
		userID, GUIDerr := profile.GetUserID(r)
		if GUIDerr != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Could not get userName")
			return
		}
		order, FUOerr := fetchUserOrder(db, f.ProductID, userID)
		if FUOerr != nil {
			log.Println(FUOerr)
			utils.RespondWithError(w, http.StatusInternalServerError, "Could not get order")
			return
		}
		json.NewEncoder(w).Encode(order)
	}
}
func PlaceOrder(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		var order models.Order
		productID := r.FormValue("productID")
		quantity := r.FormValue("quantity")
		productID32, productIDToIntErr := strconv.ParseUint(productID, 10, 32)
		if productIDToIntErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error converting productID to uint")
			return
		}
		quantity32, quantityToIntErr := strconv.ParseInt(quantity, 10, 32)
		if quantityToIntErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error converting quantity to int")
			return
		}
		userID, err := profile.GetUserID(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error getting userid")
			return
		}
		if err := db.First(&product, productID32).Error; err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("could not get the id of the specific product %d", productID32))
			return
		}

		if product.Stock < int(quantity32) {
			utils.RespondWithError(w, http.StatusInternalServerError, "error not enough to fill order")
			return
		}

		order.TotalPrice = float64(quantity32) * product.Price
		order.ProductID = uint(productID32)
		order.UserID = userID
		order.Quantity = int(quantity32)
		order.Status = "pending"
		db.Create(&order)

		orders, orderErr := fetchUserOrder(db, order.ProductID, order.UserID)
		if orderErr != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "couldn't fetch order")
			return
		}
		json.NewEncoder(w).Encode(orders)

	}
}
