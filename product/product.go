package product

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"phyEcom.com/models"
	"phyEcom.com/utils"
)

type FrontEndProductStrct struct {
	ID          uint
	Name        string
	Description string
	Price       float64
	Unique      string
	UserID      uint
	Images      []string
}

func ReorganizeForFrontend(data *[]models.Product) []FrontEndProductStrct {
	var returnData []FrontEndProductStrct
	for _, product := range *data {
		var pro FrontEndProductStrct
		pro.Name = product.Name
		pro.Description = product.Description
		pro.ID = product.ID
		pro.Unique = product.SKU
		pro.Price = product.Price
		pro.UserID = product.UserID
		for _, pic := range product.Pictures {
			pro.Images = append(pro.Images, pic.URL)
		}
		returnData = append(returnData, pro)
	}
	return returnData
}

func getPlural(word string) string { //
	if strings.HasSuffix(word, "y") {
		return word[:len(word)-1] + "ies" // e.g., "category" -> "categories"
	} else if strings.HasSuffix(word, "s") {
		return word // If already plural, return as is
	} else {
		return word + "s" // Default: add "s"
	}
}

func getBaseWord(word string) string {
	if strings.HasSuffix(word, "ies") {
		return word[:len(word)-3] + "y" // e.g., "categories" -> "category"
	} else if strings.HasSuffix(word, "es") {
		return word[:len(word)-2] // e.g., "shoes" -> "shoe"
	} else if strings.HasSuffix(word, "s") && !strings.HasSuffix(word, "ss") {
		return word[:len(word)-1] // e.g., "shirts" -> "shirt"
	}
	return word
}

func GetCount(db *gorm.DB) http.HandlerFunc {
	type count struct {
		Count int `json:"count"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var count count
		if err := db.Raw("select count(*) as count from products").Scan(&count).Error; err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "error getting product count")
			return
		}
		json.NewEncoder(w).Encode(count)
	}
}

func Product5(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var products []models.Product

		// Get pagination parameters
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		searchQuery := r.URL.Query().Get("search")

		if limit == 0 {
			limit = 10 // Default limit
		}
		offset := (page - 1) * limit

		// Preprocess search query (stemming & pluralization)
		baseQuery := getBaseWord(searchQuery) // Convert "shirts" -> "shirt"
		searchPattern := "%" + baseQuery + "%"
		pluralPattern := "%" + getPlural(baseQuery) + "%"

		// Initialize query
		query := db.Debug().Preload("Pictures")

		if searchQuery != "" {
			query = query.Where(`
				MATCH(name, description, brand) AGAINST(? IN BOOLEAN MODE)
				OR name LIKE ? OR name LIKE ?
				OR description LIKE ? OR description LIKE ?
				OR brand LIKE ? OR brand LIKE ?
				OR SOUNDEX(name) = SOUNDEX(?)
				OR SOUNDEX(description) = SOUNDEX(?)
				OR SOUNDEX(brand) = SOUNDEX(?)`,
				baseQuery,                    // Full-Text Search
				searchPattern, pluralPattern, // LIKE Matching
				searchPattern, pluralPattern, // LIKE Matching
				searchPattern, pluralPattern, // LIKE Matching
				baseQuery, baseQuery, baseQuery, // SOUNDEX for phonetic matching
			)
		}

		// Fetch products with pagination
		query.Offset(offset).Limit(limit).Order("id").Find(&products)

		data := ReorganizeForFrontend(&products)
		json.NewEncoder(w).Encode(data)
	}
}

func Product(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var products []models.Product

		// Get pagination parameters
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		searchQuery := r.URL.Query().Get("search")

		if limit == 0 {
			limit = 10 // Default limit
		}
		offset := (page - 1) * limit

		// Initialize query
		query := db.Debug().Preload("Pictures")

		if searchQuery != "" {
			// Break search query into words
			words := strings.Fields(searchQuery)
			var searchConditions []string
			var params []interface{}

			for _, word := range words {
				baseWord := getBaseWord(word) // Convert "bags" -> "bag"
				searchPattern := "%" + baseWord + "%"
				pluralPattern := "%" + getPlural(baseWord) + "%"

				// Append conditions for each word
				searchConditions = append(searchConditions, `
					MATCH(name, description, brand) AGAINST(? IN BOOLEAN MODE)
					OR name LIKE ? OR name LIKE ?
					OR description LIKE ? OR description LIKE ?
					OR brand LIKE ? OR brand LIKE ?
					OR SOUNDEX(name) = SOUNDEX(?)
					OR SOUNDEX(description) = SOUNDEX(?)
					OR SOUNDEX(brand) = SOUNDEX(?)
				`)

				params = append(params,
					baseWord,                     // Full-Text Search
					searchPattern, pluralPattern, // LIKE Matching
					searchPattern, pluralPattern, // LIKE Matching
					searchPattern, pluralPattern, // LIKE Matching
					baseWord, baseWord, baseWord, // SOUNDEX for phonetic matching
				)
			}

			// Combine conditions with OR for flexible search
			query = query.Where(strings.Join(searchConditions, " OR "), params...)
		}

		// Fetch products with pagination
		query.Offset(offset).Limit(limit).Order("id").Find(&products)

		data := ReorganizeForFrontend(&products)
		json.NewEncoder(w).Encode(data)
	}
}
