package cart

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
	"phyEcom.com/models"
	"phyEcom.com/profile"
	"phyEcom.com/utils"
)

// AddToCart handles adding an item to the cart
func AddToCart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cartItem models.Cart
		err := json.NewDecoder(r.Body).Decode(&cartItem)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "couldn't extract data")
			return
		}
		userID, uerr := profile.GetUserID(r)
		if uerr != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "user not authenticated")
			return
		}
		cartItem.UserID = userID
		var existingItem models.Cart
		result := db.Where("user_id = ? AND product_id = ?", cartItem.UserID, cartItem.ProductID).First(&existingItem)

		if result.RowsAffected > 0 {
			existingItem.Quantity += cartItem.Quantity
			db.Save(&existingItem)
			json.NewEncoder(w).Encode(existingItem)
		} else {
			db.Create(&cartItem)
			json.NewEncoder(w).Encode(cartItem)
		}
	}
}

// GetCartItems returns all cart items for a specific user
func GetCartItems(db *gorm.DB) http.HandlerFunc {
	type front struct {
		CartID          uint    `json:"cart_id"`
		ProductID       uint    `json:"product_id"`
		ProductQuantity uint    `json:"product_quantity"`
		ProductName     string  `json:"product_name"`
		ProductPrice    float64 `json:"product_price"`
		ProductPic      string  `json:"product_pic"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userID, uerr := profile.GetUserID(r)
		if uerr != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "user not authenticated")
			return
		}
		var cartItems []front
		query := `
		select c.id as cart_id, c.product_id as product_id, c.quantity as product_quantity , p.name as product_name, p.price as product_price , pi.url as product_pic 
		from carts as c
		 left join  products p on p.id = c.product_id
		 left join pictures pi on pi.id =(select id from pictures where product_id=p.id order by id limit 1)
		 where c.user_id=?
		`
		db.Raw(query, userID).Find(&cartItems)
		json.NewEncoder(w).Encode(cartItems)
	}
}

// RemoveFromCart handles removing an item from the cart
func RemoveFromCart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proid := r.URL.Query().Get("product_id")
		productIDQuery, qErr := strconv.ParseUint(proid, 10, 32)
		if qErr != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "could not parse productID")
			return
		}
		userID, uerr := profile.GetUserID(r)
		if uerr != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "user not authenticated")
			return
		}
		db.Where("user_id = ? AND product_id = ?", userID, productIDQuery).Delete(&models.Cart{})
		json.NewEncoder(w).Encode("Item removed")
	}
}

func UpdateCart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proid := r.URL.Query().Get("product_id")
		userID, uerr := profile.GetUserID(r)
		if uerr != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "user not authenticated")
			return
		}
		productIDQuery, qErr := strconv.ParseUint(proid, 10, 32)
		if qErr != nil {
			log.Println(proid)
			utils.RespondWithError(w, http.StatusBadRequest, "could not parse productID")
			return
		}
		quantity := r.URL.Query().Get("quantity")
		quantityQuery, qErr := strconv.ParseUint(quantity, 10, 32)
		if qErr != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "could not parse quantity")
			return
		}

		db.Model(&models.Cart{}).Where("user_id = ? AND product_id = ?", userID, productIDQuery).
			Update("quantity", quantityQuery)
	}
}
