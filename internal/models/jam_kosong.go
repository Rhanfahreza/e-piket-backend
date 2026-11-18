package models

import "time"

type JamKosong struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Kelas          string    `gorm:"type:varchar(50);not null" json:"kelas"`
	JamKe          string    `gorm:"type:varchar(50);not null" json:"jam_ke"`
	GuruTidakHadir string    `gorm:"type:varchar(150);not null" json:"guru_tidak_hadir"`
	GuruPengganti  string    `gorm:"type:varchar(150);not null" json:"guru_pengganti"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
