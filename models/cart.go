package models

type Cart struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
