package routes

import (
	"e-piket-be/internal/handlers"
	"e-piket-be/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	// Global Middlewares
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin, Content-Type, Accept, Authorization, X-API-Key"},
		AllowMethods: []string{"GET, POST, PUT, DELETE, OPTIONS"},
	}))

	// Base welcome route
	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"app":     "e-piket-be",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// Handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()
	jamKosongHandler := handlers.NewJamKosongHandler()

	// Health check
	app.Get("/health", healthHandler.Check)

	// API v1 routes
	api := app.Group("/api/v1")

	// 1. Public API untuk diintegrasikan di Presensi Guru (Diamankan dengan API Key: Header X-API-Key atau ?key=...)
	public := api.Group("/public")
	public.Get("/token", middleware.PresensiAPIKeyGuard(), authHandler.GetPublicDailyToken)

	// 2. Auth Routes (Login menggunakan Email + Token 4 Digit Harian)
	auth := api.Group("/auth")
	auth.Post("/login", middleware.LoginRateLimiter(), authHandler.Login)
	auth.Get("/me", middleware.JWTProtected(), authHandler.Me)

	// 3. CRUD Guru Pengganti / Jam Kosong
	jamKosong := api.Group("/jam-kosong")
	jamKosong.Get("/", jamKosongHandler.GetAll)
	jamKosong.Get("/:id", jamKosongHandler.GetByID)
	jamKosong.Post("/", jamKosongHandler.Create)
	jamKosong.Put("/:id", jamKosongHandler.Update)
	jamKosong.Delete("/:id", jamKosongHandler.Delete)

	// 4. Users Routes
	users := api.Group("/users", middleware.JWTProtected())
	users.Get("/", userHandler.GetAll)
	users.Post("/", userHandler.Create)
}
