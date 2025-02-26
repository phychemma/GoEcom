package file

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
	"phyEcom.com/auth"
	"phyEcom.com/fileUpload"
	"phyEcom.com/models"
	"phyEcom.com/utils"
)

func UploadProduct2(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		// File upload configuration

		ud := fileUpload.UploadData{
			Req:       r,
			Res:       w,
			MaxSize:   32 << 20,
			FormValue: true,
			FileName:  "image",
			FileMatch: []string{"image/jpeg", "image/png", "image/jpg"},
			Directory: "static/images/products",
			UsePath:   "/staticproductimage",
		}

		// Handle file upload
		value, paths, err := ud.NewUpload()
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to upload data")
			return
		}
		log.Println(value)
		product.Pictures = returnPicture(paths)

		product.Category = value["Categories"][0]
		product.Name = value["Name"][0]
		product.Description = value["Description"][0]
		product.Size = value["Size"][0]
		product.Color = value["Color"][0]
		product.Brand = value["Brand"][0]
		product.Material = value["Material"][0]

		userID, err := strconv.ParseUint(value["UserID"][0], 10, 32)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
			return
		}
		product.UserID = uint(userID)

		price, err := strconv.ParseFloat(value["Price"][0], 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't parse price")
			return
		}
		product.Price = price

		stock, err := strconv.ParseInt(value["Stock"][0], 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't parse stock")
			return
		}
		product.Stock = int(stock)

		// Generate SKU
		product.SKU = auth.GenerateSKU(product.Name)

		// Save product to database
		if err := db.Create(&product).Error; err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save product")
			return
		}

		// Respond with the created product
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}
