package config

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"e-piket-be/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

const tokenCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Huruf kapital & angka (tanpa O/0/I/1 yang membingungkan)

// GenerateRandomToken menghasilkan 4 karakter acak secara kriptografis (crypto/rand)
func GenerateRandomToken(length int) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(tokenCharset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = tokenCharset[num.Int64()]
	}
	return string(result), nil
}

// GetOrCreateTodayToken mengambil token hari ini, atau membuat baru jika belum ada
func GetOrCreateTodayToken(db *gorm.DB) (*models.DailyToken, error) {
	today := time.Now().Format("2006-01-02")

	var dailyToken models.DailyToken
	err := db.Where("tanggal = ?", today).First(&dailyToken).Error
	if err == nil {
		return &dailyToken, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	newTokenStr, err := GenerateRandomToken(4)
	if err != nil {
		return nil, err
	}

	dailyToken = models.DailyToken{
		Tanggal: today,
		Token:   newTokenStr,
	}

	if err := db.Create(&dailyToken).Error; err != nil {
		return nil, err
	}

	return &dailyToken, nil
}

func ConnectDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	gormLogLevel := logger.Info
	if cfg.AppEnv == "production" {
		gormLogLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke MySQL database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Berhasil terhubung ke database MySQL!")

	// Auto Migration
	err = db.AutoMigrate(
		&models.User{},
		&models.DailyToken{},
		&models.JamKosong{},
	)
	if err != nil {
		log.Printf("Peringatan: Gagal menjalankan auto-migrate: %v", err)
	} else {
		log.Println("Database AutoMigrate berhasil dijalankan.")
	}

	// Seeder Akun Default (smkpluspenus@gmail.com)
	seedDefaultUser(db)

	// Pastikan token hari ini sudah ter-generate saat server nyala
	if tokenData, err := GetOrCreateTodayToken(db); err == nil {
		log.Printf("🔑 Token Piket Hari Ini (%s): %s", tokenData.Tanggal, tokenData.Token)
	}

	DB = db
	return db
}

func seedDefaultUser(db *gorm.DB) {
	defaultEmail := "smkpluspenus@gmail.com"
	var user models.User
	if err := db.Where("email = ?", defaultEmail).First(&user).Error; err == gorm.ErrRecordNotFound {
		newUser := models.User{
			Name:  "Guru Piket SMK Plus Pelita Nusantara",
			Email: defaultEmail,
			Role:  "guru_piket",
		}
		if err := db.Create(&newUser).Error; err != nil {
			log.Printf("Gagal membuat akun default: %v", err)
		} else {
			log.Printf("✅ Akun default berhasil dibuat: %s", defaultEmail)
		}
	}
}
