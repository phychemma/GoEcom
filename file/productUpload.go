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

func returnPicture(data []string) []models.Picture {
	var pictures []models.Picture
	for _, path := range data {
		var picture models.Picture
		picture.URL = path
		pictures = append(pictures, picture)
	}
	return pictures
}

//	func returnCategory(data []string, db *gorm.DB) []*models.Category {
//		var categories []*models.Category
//		for _, path := range data {
//			var category models.Category
//			category.Name = path
//			db.Create(&category)
//			categories = append(categories, &category)
//		}
//		return categories
//	}
func UploadProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		// var categoryModal models.Category
		// var subcategoryModal models.Subcategory
		ud := fileUpload.UploadData{
			Req:       r,
			Res:       w,
			MaxSize:   32 << 20,
			FormValue: true,
			FileName:  "image",
			FileMatch: []string{"image/jpeg", "image/png", "image/jpg"},
			Directory: "static\\images\\products",
			UsePath:   "/staticproductimage",
		}
		value, paths, err := ud.NewUpload()
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to upload data")
			return
		}
		log.Print(value)
		product.Pictures = returnPicture(paths)
		// category := value["Categories"][0]
		// subcategory := value["SubCategory"][0]
		// db.Model(&models.Category{}).Select("id").Where(&models.Category{Name: category}).First(&categoryModal)
		// db.Model(&models.Subcategory{}).Select("id").Where(&models.Subcategory{Name: subcategory, CategoryID: categoryModal.ID}).First(&subcategoryModal)
		// product.SubcategoryID = subcategoryModal.ID
		product.Name = value["Name"][0]
		num64, err := strconv.ParseUint(value["UserID"][0], 10, 32)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't parse int")
			return
		}
		product.UserID = uint(num64)
		product.Description = value["Description"][0]
		flo, err := strconv.ParseFloat(value["Price"][0], 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to parse float")
			return
		}
		integ, err := strconv.ParseInt(value["Stock"][0], 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to parse int")
			return
		}
		product.Price = flo
		product.Stock = int(integ)
		product.SKU = auth.GenerateSKU(product.Name)

		db.Create(&product)
		json.NewEncoder(w).Encode(product)
	}
}
