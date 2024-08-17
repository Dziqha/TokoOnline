package models

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          string    `gorm:"column:id;primaryKey;type:VARCHAR(256)" json:"id"`
	Name        string    `gorm:"column:name;type:VARCHAR(100)" json:"name"`
	Description string    `gorm:"column:description;type:TEXT" json:"description"`
	Price       int       `gorm:"column:price" json:"price"`
	Stock       int       `gorm:"column:stock" json:"stock"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	Orders []Orders `gorm:"foreignKey:ProductId"` // One to many
}


func generateRandomId(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
    if err != nil {
        panic(err)
    }
    return hex.EncodeToString(bytes)
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
    p.ID = generateRandomId(20)
    return
}

func (p *Product) TableName() string {
	return "product"
}

func MigrateProduct(db *gorm.DB) {
	if err := db.AutoMigrate(&Product{}); err != nil {
		log.Fatalf("Error migrating products table: %v", err)
	}
}