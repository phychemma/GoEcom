package product

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"phyEcom.com/models"
	"phyEcom.com/utils"
)

type ProductOneChat struct {
	ProductID    uint
	ProductName  string
	ProductBrand string
	ProductPic   string
	LastChat     string
	ChatTime     time.Time
}

func SingleProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		productID := r.FormValue("productid")
		db.Model(&models.Product{}).Where(&models.Product{SKU: productID}).Preload("Pictures").First(&product)
		json.NewEncoder(w).Encode(ReorganizeForFrontend(&[]models.Product{product})[0])
	}
}

func GetChat(db *gorm.DB, productID uint) []models.Chat {
	var chat []models.Chat
	db.Model(&models.Chat{}).Where(&models.Chat{ProductID: productID}).Order("created_at DESC").Limit(10).Offset(0).Find(&chat)
	return chat
}
func GetProduct(db *gorm.DB, productID uint) ProductOneChat {
	var product ProductOneChat

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
            p.id = ?
    `

	if err := db.Raw(query, productID).Scan(&product).Error; err != nil {
		log.Println(err)
	}
	return product
}

func GetSingleProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		productID := r.URL.Query().Get("product_id")
		log.Println(productID)
		productIDQuery, qErr := strconv.ParseUint(productID, 10, 32)
		if qErr != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "could not parse productID")
			return
		}
		db.Debug().Preload("Pictures").Preload("carts").First(&product, productIDQuery)
		data := ReorganizeForFrontend(&[]models.Product{product})
		utils.WriteJSONResponse(w, http.StatusOK, map[string]FrontEndProductStrct{"product": data[0]})
	}
}

func ProductChat(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var productID uint64
		vars := mux.Vars(r)
		productStrVal := vars["productID"]
		if productStrVal != "" {
			var err error
			productID, err = strconv.ParseUint(productStrVal, 10, 32)
			if err != nil {
				log.Println(err)
				utils.RespondWithError(w, http.StatusBadRequest, "error coverting productID to an integer")
				return
			}
		}
		product := GetProduct(db, uint(productID))
		chat := GetChat(db, uint(productID))
		utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"product": product, "chat": chat})
	}
}
