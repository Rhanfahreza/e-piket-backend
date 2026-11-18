package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Role      string         `gorm:"type:varchar(50);default:'guru_piket'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type DailyToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Tanggal   string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"tanggal"` // Format: YYYY-MM-DD
	Token     string    `gorm:"type:varchar(10);not null" json:"token"`               // 4 Karakter Huruf & Angka
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
