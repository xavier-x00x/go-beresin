package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"go-beresin/internal/repository"
	"go-beresin/internal/transport/middleware"
	"go-beresin/pkg/utils"
)

// AuthHandler orchestrates authentication and security actions.
type AuthHandler struct {
	q       *repository.Queries
	db      *pgxpool.Pool
	rdb     *redis.Client
	limiter *middleware.RateLimiter
}

// NewAuthHandler constructs a new AuthHandler.
func NewAuthHandler(db *pgxpool.Pool, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{
		q:       repository.New(db),
		db:      db,
		rdb:     rdb,
		limiter: middleware.NewRateLimiter(rdb),
	}
}

// ---------------------------------------------------------
// REQUEST REQUEST STRUCTS
// ---------------------------------------------------------

type RegisterReq struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"` // 'user', 'talent'
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GoogleLoginReq struct {
	Token string `json:"token"`
}

type SendOtpReq struct {
	Phone string `json:"phone"`
}

type VerifyOtpReq struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

type ResetPasswordReq struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// Helper: Convert pgtype.UUID to standard string representation
func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	var src [16]byte
	copy(src[:], u.Bytes[:])
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

// Helper: Convert string to pgtype.UUID
func stringToUUID(s string) pgtype.UUID {
	var uuid pgtype.UUID
	_ = uuid.Scan(s)
	return uuid
}

// Helper: Validate Register inputs
func validateRegisterInput(email, phone, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email format")
	}
	phoneRegex := regexp.MustCompile(`^(?:\+62|62|0)8[1-9][0-9]{6,10}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid Indonesian phone number format")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

// ---------------------------------------------------------
// ENDPOINTS IMPLEMENTATIONS
// ---------------------------------------------------------

// Register handles user registration with email, phone, and password.
// @Summary User Registration
// @Description Register a new client or talent account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterReq true "Register Payload"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 409 {object} Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(RegisterReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	// Default role to 'user' if empty
	role := strings.ToLower(req.Role)
	if role == "" {
		role = "user"
	}
	if role != "user" && role != "talent" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Role must be 'user' or 'talent'",
		})
	}

	// Validate inputs
	if err := validateRegisterInput(req.Email, req.Phone, req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	ctx := context.Background()

	// Check email duplicate
	_, err := h.q.GetUserByEmail(ctx, pgtype.Text{String: req.Email, Valid: true})
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(Response{
			Status:  "error",
			Message: "Email is already registered",
		})
	}

	// Check phone duplicate
	_, err = h.q.GetUserByPhone(ctx, req.Phone)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(Response{
			Status:  "error",
			Message: "Phone number is already registered",
		})
	}

	// Hash password (salt/cost 12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to securely process password",
		})
	}

	// Save user to DB
	user, err := h.q.CreateUser(ctx, repository.CreateUserParams{
		Email:        pgtype.Text{String: req.Email, Valid: true},
		Phone:        req.Phone,
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
		FullName:     req.FullName,
		Role:         role,
		IsVerified:   pgtype.Bool{Bool: false, Valid: true},
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to register user to database",
		})
	}

	userIDStr := uuidToString(user.ID)

	// Generate verification token (Simulated Job)
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	verifyToken := hex.EncodeToString(b)

	// Save verification token in Redis (TTL 24 hours)
	h.rdb.Set(ctx, fmt.Sprintf("auth_verify_token:%s", verifyToken), userIDStr, 24*time.Hour)

	log.Printf("[BACKGROUND JOB] Enqueued verification email with token %s to %s", verifyToken, req.Email)

	// Write Audit Log
	_, _ = h.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		UserID:    user.ID,
		Action:    "register",
		IpAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	})

	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:  "success",
		Message: "Registration completed successfully. Please check your email to verify your account.",
		Data: fiber.Map{
			"user_id":   userIDStr,
			"email":     user.Email.String,
			"phone":     user.Phone,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}

// Login verifies credentials and generates session tokens (access + refresh).
// @Summary User Login
// @Description Login with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginReq true "Login Payload"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 429 {object} Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ip := c.IP()
	ctx := context.Background()

	// Check failed login rate-limiter blocking state
	blocked, ttl, err := h.limiter.IsIPBlocked(ctx, ip)
	if err != nil {
		log.Printf("[WARNING] Rate limiter error: %v", err)
	}
	if blocked {
		return c.Status(fiber.StatusTooManyRequests).JSON(Response{
			Status:  "error",
			Message: fmt.Sprintf("IP temporarily blocked due to multiple failed login attempts. Please retry after %d seconds.", int(ttl.Seconds())),
		})
	}

	req := new(LoginReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Email and password are required",
		})
	}

	// Fetch user from DB
	user, err := h.q.GetUserByEmail(ctx, pgtype.Text{String: req.Email, Valid: true})
	if err != nil {
		// Increment failed counter
		_, _ = h.limiter.RecordFailedLogin(ctx, ip)
		return c.Status(fiber.StatusUnauthorized).JSON(Response{
			Status:  "error",
			Message: "Invalid email or password credentials",
		})
	}

	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password))
	if err != nil {
		// Increment failed counter
		_, _ = h.limiter.RecordFailedLogin(ctx, ip)
		return c.Status(fiber.StatusUnauthorized).JSON(Response{
			Status:  "error",
			Message: "Invalid email or password credentials",
		})
	}

	// Success! Reset failed counter
	_ = h.limiter.ResetFailedLogin(ctx, ip)

	userIDStr := uuidToString(user.ID)

	// Generate Access (15m) and Refresh (30d) token pair
	accessToken, refreshToken, err := utils.GenerateTokenPair(userIDStr, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to generate session tokens",
		})
	}

	// Store Refresh Token in Redis (TTL 30 days)
	err = h.rdb.Set(ctx, fmt.Sprintf("user:%s:refresh", userIDStr), refreshToken, 30*24*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to register session",
		})
	}

	// Write Audit Log
	_, _ = h.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		UserID:    user.ID,
		Action:    "login",
		IpAddress: ip,
		UserAgent: c.Get("User-Agent"),
	})

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully authenticated",
		Data: fiber.Map{
			"user_id":       userIDStr,
			"full_name":     user.FullName,
			"role":          user.Role,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// GoogleLogin handles Google OAuth2 verification and logs user in.
// @Summary Google Social Login
// @Description Login or register automatically via Google OAuth token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body GoogleLoginReq true "Google Token Payload"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/auth/login/google [post]
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	req := new(GoogleLoginReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := context.Background()

	// Verify Google Token
	gUser, err := utils.VerifyGoogleToken(ctx, req.Token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	var user repository.User

	// Check if google_id already exists in DB
	user, err = h.q.GetUserByGoogleID(ctx, pgtype.Text{String: gUser.ID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Google ID not found. Check if email exists
			existingUser, emailErr := h.q.GetUserByEmail(ctx, pgtype.Text{String: gUser.Email, Valid: true})
			if emailErr == nil {
				// Email exists. Link Google account!
				user, err = h.q.UpdateUserGoogleID(ctx, repository.UpdateUserGoogleIDParams{
					ID:       existingUser.ID,
					GoogleID: pgtype.Text{String: gUser.ID, Valid: true},
				})
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(Response{
						Status:  "error",
						Message: "Failed to link Google account",
					})
				}
				// Set verified status
				_, _ = h.q.UpdateUserVerificationStatus(ctx, repository.UpdateUserVerificationStatusParams{
					ID:         user.ID,
					IsVerified: pgtype.Bool{Bool: true, Valid: true},
				})
			} else {
				// User email doesn't exist either. Create new Google user account.
				// Since phone column is UNIQUE & NOT NULL, we create a deterministic mock phone
				uniqueMockPhone := "+62899" + gUser.ID[len(gUser.ID)-8:]
				user, err = h.q.CreateUser(ctx, repository.CreateUserParams{
					Email:      pgtype.Text{String: gUser.Email, Valid: true},
					Phone:      uniqueMockPhone,
					FullName:   gUser.Name,
					Role:       "user",
					AvatarUrl:  pgtype.Text{String: gUser.AvatarURL, Valid: true},
					IsVerified: pgtype.Bool{Bool: true, Valid: true},
					IsActive:   pgtype.Bool{Bool: true, Valid: true},
					GoogleID:   pgtype.Text{String: gUser.ID, Valid: true},
				})
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(Response{
						Status:  "error",
						Message: "Failed to register new Google account",
					})
				}
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(Response{
				Status:  "error",
				Message: "Failed to query database user",
			})
		}
	}

	userIDStr := uuidToString(user.ID)

	// Generate Access and Refresh token pair
	accessToken, refreshToken, err := utils.GenerateTokenPair(userIDStr, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to generate session tokens",
		})
	}

	// Store Refresh Token in Redis
	h.rdb.Set(ctx, fmt.Sprintf("user:%s:refresh", userIDStr), refreshToken, 30*24*time.Hour)

	// Write Audit Log
	_, _ = h.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		UserID:    user.ID,
		Action:    "login_google",
		IpAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	})

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully authenticated via Google",
		Data: fiber.Map{
			"user_id":       userIDStr,
			"full_name":     user.FullName,
			"role":          user.Role,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// SendWhatsAppOTP generates a 6-digit OTP code and dispatches it via Twilio.
// @Summary Send WhatsApp OTP
// @Description Request a 6-digit WhatsApp login OTP
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SendOtpReq true "Send OTP Payload"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/auth/login/whatsapp/send-otp [post]
func (h *AuthHandler) SendWhatsAppOTP(c *fiber.Ctx) error {
	req := new(SendOtpReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	phoneRegex := regexp.MustCompile(`^(?:\+62|62|0)8[1-9][0-9]{6,10}$`)
	if !phoneRegex.MatchString(req.Phone) {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid Indonesian phone number format",
		})
	}

	ctx := context.Background()

	// Generate 6-digit OTP
	otp := utils.GenerateOTP()

	// Save OTP in Redis (TTL 5 minutes)
	err := h.rdb.Set(ctx, fmt.Sprintf("whatsapp_otp:%s", req.Phone), otp, 5*time.Minute).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to process OTP request",
		})
	}

	// Dispatch OTP via WhatsApp helper
	err = utils.SendWhatsAppOTP(ctx, req.Phone, otp)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to dispatch WhatsApp OTP message",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "OTP has been successfully dispatched to WhatsApp",
	})
}

// VerifyWhatsAppOTP checks OTP from Redis and logs the user in.
// @Summary Verify WhatsApp OTP
// @Description Verify WhatsApp OTP code and start login session
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body VerifyOtpReq true "Verify OTP Payload"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/auth/login/whatsapp/verify [post]
func (h *AuthHandler) VerifyWhatsAppOTP(c *fiber.Ctx) error {
	req := new(VerifyOtpReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	if req.Phone == "" || req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Phone number and OTP are required",
		})
	}

	ctx := context.Background()
	otpKey := fmt.Sprintf("whatsapp_otp:%s", req.Phone)

	// Fetch stored OTP
	storedOtp, err := h.rdb.Get(ctx, otpKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return c.Status(fiber.StatusBadRequest).JSON(Response{
				Status:  "error",
				Message: "OTP has expired or does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to verify OTP",
		})
	}

	// Verify OTP match
	if storedOtp != req.OTP {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Incorrect OTP code",
		})
	}

	// OTP matched successfully. Delete OTP (one-time use)
	h.rdb.Del(ctx, otpKey)

	var user repository.User

	// Check if user exists with this phone
	user, err = h.q.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// User phone not found. Create new verified account automatically.
			user, err = h.q.CreateUser(ctx, repository.CreateUserParams{
				Phone:      req.Phone,
				FullName:   "WhatsApp User",
				Role:       "user",
				IsVerified: pgtype.Bool{Bool: true, Valid: true},
				IsActive:   pgtype.Bool{Bool: true, Valid: true},
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(Response{
					Status:  "error",
					Message: "Failed to automatically register new account",
				})
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(Response{
				Status:  "error",
				Message: "Failed to check account state",
			})
		}
	}

	userIDStr := uuidToString(user.ID)

	// Generate Access and Refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokenPair(userIDStr, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to generate session tokens",
		})
	}

	// Store Refresh Token in Redis
	h.rdb.Set(ctx, fmt.Sprintf("user:%s:refresh", userIDStr), refreshToken, 30*24*time.Hour)

	// Write Audit Log
	_, _ = h.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		UserID:    user.ID,
		Action:    "login_whatsapp",
		IpAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	})

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully authenticated via WhatsApp OTP",
		Data: fiber.Map{
			"user_id":       userIDStr,
			"full_name":     user.FullName,
			"role":          user.Role,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// RefreshToken validates the refresh token from Redis and rotates it.
// @Summary Refresh Token Rotation
// @Description Rotate expired access token using valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenReq true "Refresh Token Payload"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /api/v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	req := new(RefreshTokenReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Refresh token is required",
		})
	}

	// Validate token signature and expiry
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(Response{
			Status:  "error",
			Message: "Invalid or expired refresh token",
		})
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("user:%s:refresh", claims.UserID)

	// Fetch current valid Refresh Token from Redis
	storedToken, err := h.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(Response{
			Status:  "error",
			Message: "Session has expired or logout completed. Please login again.",
		})
	}

	// Verify token match (detect refresh token reuse/hijack)
	if storedToken != req.RefreshToken {
		// Revoke session completely as warning measure
		h.rdb.Del(ctx, redisKey)
		return c.Status(fiber.StatusUnauthorized).JSON(Response{
			Status:  "error",
			Message: "Security warning: refresh token mismatch. Session revoked.",
		})
	}

	// Rotate: generate new Access & Refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokenPair(claims.UserID, claims.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to rotate session tokens",
		})
	}

	// Update Redis with new rotated Refresh Token (TTL 30 days)
	h.rdb.Set(ctx, redisKey, refreshToken, 30*24*time.Hour)

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Tokens successfully rotated",
		Data: fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// Logout deletes the refresh token from Redis.
// @Summary User Logout
// @Description Invalidate user active session and revoke refresh token
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(Response{
			Status:  "error",
			Message: "Missing authenticated user details",
		})
	}

	userIDStr := userID.(string)
	ctx := context.Background()

	// Delete Refresh token from Redis
	h.rdb.Del(ctx, fmt.Sprintf("user:%s:refresh", userIDStr))

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully logged out, session terminated",
	})
}

// ForgotPassword generates reset token, saves in Redis, and simulates email.
// @Summary Forgot Password
// @Description Request a password reset email token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ForgotPasswordReq true "Forgot Password Payload"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	req := new(ForgotPasswordReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Email address is required",
		})
	}

	ctx := context.Background()

	// Check if user exists
	user, err := h.q.GetUserByEmail(ctx, pgtype.Text{String: req.Email, Valid: true})
	if err != nil {
		// Fail silently to prevent user email enumeration security issue
		return c.Status(fiber.StatusOK).JSON(Response{
			Status:  "success",
			Message: "If the email exists, password reset instructions have been dispatched.",
		})
	}

	// Generate secure random reset token
	b := make([]byte, 20)
	_, _ = rand.Read(b)
	resetToken := hex.EncodeToString(b)

	// Save reset token mapping to user ID in Redis (TTL 1 hour)
	userIDStr := uuidToString(user.ID)
	h.rdb.Set(ctx, fmt.Sprintf("reset_password_token:%s", resetToken), userIDStr, 1*time.Hour)

	log.Printf("[BACKGROUND JOB] Enqueued password reset email with token %s to %s", resetToken, req.Email)

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "If the email exists, password reset instructions have been dispatched.",
	})
}

// ResetPassword validates the reset token and updates password hash.
// @Summary Reset Password
// @Description Reset password utilizing valid email reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResetPasswordReq true "Reset Password Payload"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	req := new(ResetPasswordReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	if req.Token == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Reset token and new password are required",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "New password must be at least 8 characters long",
		})
	}

	ctx := context.Background()
	tokenKey := fmt.Sprintf("reset_password_token:%s", req.Token)

	// Get User ID from Redis
	userIDStr, err := h.rdb.Get(ctx, tokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return c.Status(fiber.StatusBadRequest).JSON(Response{
				Status:  "error",
				Message: "Reset token is invalid or has expired",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to verify reset token",
		})
	}

	// Encrypt new password (salt/cost 12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to secure new password",
		})
	}

	// Update DB password
	userID := stringToUUID(userIDStr)
	_, err = h.q.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to update password in database",
		})
	}

	// Security Measure: Invalidate all active sessions for this user in Redis
	h.rdb.Del(ctx, fmt.Sprintf("user:%s:refresh", userIDStr))

	// Delete used reset token
	h.rdb.Del(ctx, tokenKey)

	// Write Audit Log
	_, _ = h.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		UserID:    userID,
		Action:    "reset_password",
		IpAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	})

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Password updated successfully. All other sessions have been logged out.",
	})
}
