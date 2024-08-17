package models

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"gorm.io/gorm"
)

type Carts struct {
	gorm.Model
	ID        string         `gorm:"column:id;primaryKey;type:VARCHAR(256)" json:"id"`
	Quantity  int            `gorm:"column:quantity" json:"quantity" validate:"required"`
	UserId    string         `gorm:"column:user_id;type:VARCHAR(256)" json:"userId" validate:"required"`
	ProductId string         `gorm:"column:product_id;type:VARCHAR(256)" json:"productId" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	User    Users   `gorm:"foreignKey:UserId"`    // Relasi ke Users
	Product Product `gorm:"foreignKey:ProductId"` // Relasi ke Product
}



func generateRandomIdCarts(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func (c *Carts) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = generateRandomIdCarts(20)
	return
}

func (c * Carts) TableName() string {
	return "carts"
}


func MigrateCarts(db *gorm.DB) {
	if err := db.AutoMigrate(&Carts{}); err != nil {
		log.Fatalf("Error migrating carts table: %v", err)
	}
}