package models

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID       string    `gorm:"column:id;primaryKey;type:VARCHAR(256)" json:"id"`
	Username string    `gorm:"column:username;type:VARCHAR(100)" json:"username" validate:"required"`
	Password string    `gorm:"column:password;type:VARCHAR(100)" json:"password" validate:"required"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	Orders []Orders `gorm:"foreignKey:UserId"` // One to many
}


func generateRandomIdUser(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
    if err != nil {
        panic(err)
    }
    return hex.EncodeToString(bytes)
}

func (u *Users) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = generateRandomIdUser(20)
    return
}


func (s *Users) TableName() string {
	return "users"
}

func MigrateUsers(db *gorm.DB) {
	if err := db.AutoMigrate(&Carts{}); err != nil {
		log.Fatalf("Error migrating carts table: %v", err)
	}
}
