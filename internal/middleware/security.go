package middleware

import (
	"crypto/subtle"
	"strings"
	"sync"
	"time"

	"e-piket-be/config"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// ==========================================
// 1. RATE LIMITER KHUSUS LOGIN (ANTI BRUTE-FORCE)
// ==========================================

type ipAttempt struct {
	count     int
	expiresAt time.Time
}

var (
	loginAttempts = make(map[string]*ipAttempt)
	mu            sync.Mutex
)

const (
	maxLoginAttempts = 5
	blockDuration    = 15 * time.Minute
)

// LoginRateLimiter membatasi percobaan login maksimal 5 kali gagal per 15 menit per IP
func LoginRateLimiter() fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()
		now := time.Now()

		mu.Lock()
		attempt, exists := loginAttempts[ip]
		if exists {
			if now.After(attempt.expiresAt) {
				delete(loginAttempts, ip)
				exists = false
			} else if attempt.count >= maxLoginAttempts {
				sisaDetik := int(time.Until(attempt.expiresAt).Seconds())
				mu.Unlock()
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success":             false,
					"message":             "Terlalu banyak percobaan login gagal. Silakan tunggu beberapa saat sebelum mencoba lagi.",
					"retry_after_seconds": sisaDetik,
				})
			}
		}
		mu.Unlock()

		return c.Next()
	}
}

// RecordFailedLogin mencatat kegagalan login untuk IP tertentu
func RecordFailedLogin(ip string) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	attempt, exists := loginAttempts[ip]
	if !exists || now.After(attempt.expiresAt) {
		loginAttempts[ip] = &ipAttempt{
			count:     1,
			expiresAt: now.Add(blockDuration),
		}
		return
	}
	attempt.count++
}

// ResetLoginAttempts menghapus catatan gagal jika user berhasil login
func ResetLoginAttempts(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(loginAttempts, ip)
}

// ==========================================
// 2. API KEY GUARD UNTUK ENDPOINT PUBLIC TOKEN
// ==========================================

// PresensiAPIKeyGuard mengamankan endpoint pengambilan token harian
// Mendukung 2 cara: Header "X-API-Key" atau Query Parameter "?key="
func PresensiAPIKeyGuard() fiber.Handler {
	return func(c fiber.Ctx) error {
		expectedKey := "penus2026secret"
		if config.AppConfig != nil && config.AppConfig.PresensiAPIKey != "" {
			expectedKey = config.AppConfig.PresensiAPIKey
		}

		// Cek dari Header X-API-Key terlebih dahulu, jika kosong cek dari query ?key=
		providedKey := strings.TrimSpace(c.Get("X-API-Key"))
		if providedKey == "" {
			providedKey = strings.TrimSpace(c.Query("key"))
		}

		if providedKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak: API Key diperlukan (gunakan header X-API-Key atau parameter ?key=)",
			})
		}

		// Constant-time comparison untuk mencegah Timing Attack
		if subtle.ConstantTimeCompare([]byte(providedKey), []byte(expectedKey)) != 1 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak: API Key tidak valid",
			})
		}

		return c.Next()
	}
}

// ==========================================
// 3. JWT AUTH MIDDLEWARE (SESI HARIAN)
// ==========================================

func JWTProtected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Akses ditolak: Token sesi tidak ditemukan",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := "epiket_smkpluspenus_super_secret_jwt_2026"
		if config.AppConfig != nil && config.AppConfig.JWTSecret != "" {
			secret = config.AppConfig.JWTSecret
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Sesi tidak valid atau sudah berakhir. Silakan login kembali dengan token hari ini.",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Klaim token tidak valid",
			})
		}

		// Pastikan token JWT berasal dari tanggal hari ini
		today := time.Now().Format("2006-01-02")
		if loginDate, ok := claims["login_date"].(string); !ok || loginDate != today {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Sesi hari sebelumnya telah berakhir. Silakan login menggunakan token hari ini.",
			})
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("email", claims["email"])
		c.Locals("role", claims["role"])
		c.Locals("login_date", claims["login_date"])

		return c.Next()
	}
}
