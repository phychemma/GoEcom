package product

// import (
// 	"encoding/json"
// 	"net/http"

// 	"gorm.io/gorm"
// 	"phyEcom.com/models"
// )

// func FetchCategory(db *gorm.DB) http.HandlerFunc {

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var category []models.Category
// 		db.Model(&models.Category{}).Select([]string{"name", "id"}).Find(&category)
// 		json.NewEncoder(w).Encode(category)
// 	}
// }

// func FetchSubcategory(db *gorm.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var category models.Category
// 		var subcategory []models.Subcategory
// 		sentCategory := r.FormValue("category")
// 		db.Model(&models.Category{}).Where(&models.Category{Name: sentCategory}).First(&category)
// 		db.Model(&models.Subcategory{}).Where(&models.Subcategory{CategoryID: category.ID}).Find(&subcategory)
// 		json.NewEncoder(w).Encode(subcategory)

// 	}
// }
