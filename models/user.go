package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID               uint `gorm:"primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Username         string `gorm:"unique"`
	Password         string
	Email            string `gorm:"unique"`
	VerificationCode string
	EmailVerified    bool     `gorm:"default:false"`
	Reviews          []Review `gorm:"foreignKey:UserID"`
	Role             string
	Profile          Profile   `gorm:"constraint:OnDelete:CASCADE;"`
	Product          []Product `gorm:"forignKey:UserID"`
	Chats            []Chat    `gorm:"foreignKey:UserID"` // List of chats initiated by the user

}

type Profile struct {
	ID        uint `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Image     string
	UserID    uint `gorm:"unique;not null"`
}

type UserExist struct {
	Username string
	Exist    bool
}

type SessionData struct {
}

type EmailExist struct {
	Email string
	Exist bool
}

type JwtToken struct {
	Token string `json:"token"`
}
type LoggedIn struct {
	Loggedin bool
}
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
type VerifyCode struct {
	Code  string
	Email string
}

type Product struct {
	UserID      uint
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text;not null"`
	Category    string
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	Stock       int       `gorm:"not null"`
	SKU         string    `gorm:"type:varchar(100);unique"`
	Pictures    []Picture `gorm:"foreignKey:ProductID"`
	Reviews     []Review  `gorm:"foreignKey:ProductID"`
	Cart        Cart      `gorm:"foreignKey:ProductID"`
	Chats       []Chat    `gorm:"foreignKey:ProductID"`       // Chats related to this product
	Size        string    `gorm:"type:varchar(50);not null"`  // Size of the wear (e.g., S, M, L, XL)
	Color       string    `gorm:"type:varchar(50);not null"`  // Color of the wear
	Brand       string    `gorm:"type:varchar(100);not null"` // Brand or designer
	Material    string    `gorm:"type:varchar(100)"`          // Material (e.g., Cotton, Polyester)
}

// Order model
type Order struct {
	ID         uint    `gorm:"primaryKey"`
	UserID     uint    `gorm:"not null"`
	ProductID  uint    `gorm:"not null"`
	Quantity   int     `gorm:"not null"`
	TotalPrice float64 `gorm:"not null"`
	Status     string  `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       User    `gorm:"foreignKey:UserID"`
	Product    Product `gorm:"foreignKey:ProductID"`
	Seller     bool    `gorm:"default:false"`
	Buyer      bool    `gorm:"default:false"`
}

type OrderDelivered struct {
	ID      uint `gorm:"not null"`
	OrderID uint `gorm:"not null"`
	Seller  bool `gorm:"default:false"`
	Buyer   bool `gorm:"default:false"`
}

type Picture struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	URL       string `gorm:"type:varchar(255);not null"`
	ProductID uint   `gorm:"not null"`
	//Product   Product `gorm:"foreignKey:ProductID"`
}
type Review struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Rating    int    `gorm:"type:int;not null"`
	Comment   string `gorm:"type:text"`
	ProductID uint   `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	User      *User  `gorm:"foreignKey:UserID"`
}
