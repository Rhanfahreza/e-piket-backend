package handlers

import (
	"crypto/subtle"
	"strings"
	"time"

	"e-piket-be/config"
	"e-piket-be/internal/middleware"
	"e-piket-be/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type LoginInput struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

// GetPublicDailyToken mengembalikan token 4 digit aktif hari ini untuk diintegrasikan ke Presensi Guru
func (h *AuthHandler) GetPublicDailyToken(c fiber.Ctx) error {
	dailyToken, err := config.GetOrCreateTodayToken(config.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyiapkan token harian",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Token harian piket aktif",
		"data": fiber.Map{
			"tanggal":      dailyToken.Tanggal,
			"token":        dailyToken.Token,
			"berlaku_akun": "smkpluspenus@gmail.com",
		},
	})
}

// Login memverifikasi email smkpluspenus@gmail.com + token 4 digit hari ini
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var input LoginInput
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format JSON tidak valid",
		})
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	inputToken := strings.ToUpper(strings.TrimSpace(input.Token))

	if email == "" || inputToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email dan token 4 digit wajib diisi",
		})
	}

	if len(inputToken) != 4 {
		middleware.RecordFailedLogin(c.IP())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Email atau token harian tidak valid",
		})
	}

	// 1. Cari user berdasarkan email
	var user models.User
	userErr := config.DB.Where("email = ?", email).First(&user).Error

	// 2. Ambil token aktif hari ini
	dailyToken, tokenErr := config.GetOrCreateTodayToken(config.DB)
	if tokenErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Terjadi kesalahan sistem saat memverifikasi token",
		})
	}

	// 3. Bandingkan token secara Constant-Time (mencegah Timing Attack)
	tokenMatch := subtle.ConstantTimeCompare([]byte(inputToken), []byte(dailyToken.Token)) == 1

	if userErr != nil || !tokenMatch {
		middleware.RecordFailedLogin(c.IP())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Email atau token harian salah",
		})
	}

	// Reset hitungan gagal login karena sudah berhasil
	middleware.ResetLoginAttempts(c.IP())

	// 4. Buat JWT yang otomatis hangus pada akhir hari ini (23:59:59)
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	secret := "epiket_smkpluspenus_super_secret_jwt_2026"
	if config.AppConfig != nil && config.AppConfig.JWTSecret != "" {
		secret = config.AppConfig.JWTSecret
	}

	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"login_date": dailyToken.Tanggal,
		"iat":        now.Unix(),
		"exp":        endOfDay.Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat sesi login",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data": fiber.Map{
			"access_token": signedToken,
			"tanggal":      dailyToken.Tanggal,
			"expires_at":   endOfDay.Format("2006-01-02 15:04:05"),
			"user": fiber.Map{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		},
	})
}

// Me memeriksa profil dan sesi aktif dari JWT token
func (h *AuthHandler) Me(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id":    c.Locals("user_id"),
			"email":      c.Locals("email"),
			"role":       c.Locals("role"),
			"login_date": c.Locals("login_date"),
		},
	})
}
