package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go-beresin/internal/service"
	"go-beresin/internal/transport/handler"
	"go-beresin/internal/transport/middleware"
)

// SetupRoutes registers all routes and middlewares.
func SetupRoutes(app *fiber.App, rdb *redis.Client, dbPool *pgxpool.Pool) {
	// 1. Global Middlewares
	app.Use(recover.New()) // Capture panic and prevent crash
	app.Use(cors.New())    // Enable Cross-Origin Resource Sharing
	app.Use(logger.New(logger.Config{
		Format: "[HTTP] ${time} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
	}))

	// 2. Global Rate Limiter (60 requests per minute)
	limiter := middleware.NewRateLimiter(rdb)
	app.Use(limiter.GlobalLimit())

	// 3. Swagger Route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// 4. Initialise Service & Handlers
	h := handler.NewDummyHandler()
	authSvc := service.NewAuthService(dbPool, rdb)
	authH := handler.NewAuthHandler(authSvc)

	// 5. System Health Check (exempt from JWT)
	app.Get("/health", h.HealthCheck)

	// 6. API v1 Router Group
	v1 := app.Group("/api/v1")

	// Strict rate limiting specifically for auth endpoints (5 requests per minute)
	auth := v1.Group("/auth", limiter.StrictLimit())
	auth.Post("/register", authH.Register)
	auth.Post("/login", authH.Login)
	auth.Post("/login/google", authH.GoogleLogin)
	auth.Post("/login/whatsapp/send-otp", authH.SendWhatsAppOTP)
	auth.Post("/login/whatsapp/verify", authH.VerifyWhatsAppOTP)
	auth.Post("/refresh-token", authH.RefreshToken)
	auth.Post("/logout", middleware.JWTAuth(), authH.Logout)
	auth.Post("/forgot-password", authH.ForgotPassword)
	auth.Post("/reset-password", authH.ResetPassword)

	// Protected Routes (requires valid JWT token)
	protected := v1.Group("/", middleware.JWTAuth())
	
	// User profile
	protected.Get("/profile", h.Profile)
	
	// File Upload Handler (Multipart)
	protected.Post("/upload", h.Upload)

	// Admin Only Dashboard (requires JWT and 'admin' role)
	adminOnly := v1.Group("/admin", middleware.JWTAuth(), middleware.RequireRole("admin"))
	adminOnly.Get("/", h.AdminOnly)
}
