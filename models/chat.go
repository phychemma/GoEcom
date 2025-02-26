package models

import (
	"time"
)

// Chat represents a user-initiated conversation related to a product
type Chat struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     uint    `gorm:"not null"`                     // Foreign key for User
	ProductID  uint    `gorm:"not null"`                     // Foreign key for Product
	AdminID    uint    `gorm:"default:0"`                    // ID for the admin responding to the chat
	Message    string  `gorm:"type:text;not null"`           // The chat message content
	IsAdmin    bool    `gorm:"default:false"`                // Whether the message is from an admin
	Read       bool    `gorm:"default:false"`                // Whether the message has been read
	Attachment string  `gorm:"type:varchar(255)"`            // Optional URL for an attachment
	Product    Product `gorm:"constraint:OnDelete:CASCADE;"` // Relationship with Product
	User       User    `gorm:"constraint:OnDelete:CASCADE;"` // Relationship with User
}
