package models

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"gorm.io/gorm"
)

type Orders struct {
	gorm.Model
	ID         string `gorm:"column:id;primaryKey;type:VARCHAR(256)" json:"id"`
	UserId     string `gorm:"column:user_id;type:VARCHAR(256)" json:"userId" validate:"required"`
	ProductId  string `gorm:"column:product_id;type:VARCHAR(256)" json:"productId" validate:"required"`
	Quantity   int    `gorm:"column:quantity" json:"quantity" validate:"required"`
	TotalPrice int    `gorm:"column:total_price" json:"total_price" validate:"required"`
	Status     string `gorm:"column:status;type:VARCHAR(256);default:'pending'" json:"status"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	User     Users   `gorm:"foreignKey:UserId"`   // Relasi ke Users
	Product  Product `gorm:"foreignKey:ProductId"` // Relasi ke Product
}


func generateRandomIdOrder(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
    if err != nil {
        panic(err)
    }
    return hex.EncodeToString(bytes)
}

func (o *Orders) BeforeCreate(tx *gorm.DB) (err error) {
    o.ID = generateRandomIdOrder(20)
    return
}

func (o *Orders) TableName() string {
	return "orders"
}

func MigrateOrder(db *gorm.DB) {
	if err := db.AutoMigrate(&Orders{}); err != nil {
		log.Fatalf("Error migrating orders table: %v", err)
	}
}