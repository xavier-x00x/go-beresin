package service

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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"go-beresin/internal/domain"
	"go-beresin/internal/repository"
	"go-beresin/internal/transport/middleware"
	"go-beresin/pkg/utils"
)

type AuthService interface {
	Register(ctx context.Context, req domain.RegisterReq, ip, userAgent string) (*domain.RegisterResp, error)
	Login(ctx context.Context, req domain.LoginReq, ip, userAgent string) (*domain.LoginResp, error)
	GoogleLogin(ctx context.Context, req domain.GoogleLoginReq, ip, userAgent string) (*domain.LoginResp, error)
	SendWhatsAppOTP(ctx context.Context, req domain.SendOtpReq) error
	VerifyWhatsAppOTP(ctx context.Context, req domain.VerifyOtpReq, ip, userAgent string) (*domain.LoginResp, error)
	RefreshToken(ctx context.Context, req domain.RefreshTokenReq) (*domain.RefreshTokenResp, error)
	Logout(ctx context.Context, userID string) error
	ForgotPassword(ctx context.Context, req domain.ForgotPasswordReq) error
	ResetPassword(ctx context.Context, req domain.ResetPasswordReq, ip, userAgent string) error
}

type authService struct {
	q       *repository.Queries
	db      *pgxpool.Pool
	rdb     *redis.Client
	limiter *middleware.RateLimiter
}

func NewAuthService(db *pgxpool.Pool, rdb *redis.Client) AuthService {
	return &authService{
		q:       repository.New(db),
		db:      db,
		rdb:     rdb,
		limiter: middleware.NewRateLimiter(rdb),
	}
}

func validateRegisterInput(email, phone, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return domain.NewValidationError("invalid email format")
	}
	phoneRegex := regexp.MustCompile(`^(?:\+62|62|0)8[1-9][0-9]{6,10}$`)
	if !phoneRegex.MatchString(phone) {
		return domain.NewValidationError("invalid Indonesian phone number format")
	}
	if len(password) < 8 {
		return domain.NewValidationError("password must be at least 8 characters long")
	}
	return nil
}

func phoneRegex() *regexp.Regexp {
	return regexp.MustCompile(`^(?:\+62|62|0)8[1-9][0-9]{6,10}$`)
}

func (s *authService) Register(ctx context.Context, req domain.RegisterReq, ip, userAgent string) (*domain.RegisterResp, error) {
	role := strings.ToLower(req.Role)
	if role == "" {
		role = "user"
	}
	if role != "user" && role != "talent" {
		return nil, domain.NewValidationError("role must be 'user' or 'talent'")
	}

	if err := validateRegisterInput(req.Email, req.Phone, req.Password); err != nil {
		return nil, err
	}

	_, err := s.q.GetUserByEmail(ctx, pgtype.Text{String: req.Email, Valid: true})
	if err == nil {
		return nil, domain.ErrEmailConflict
	}

	_, err = s.q.GetUserByPhone(ctx, req.Phone)
	if err == nil {
		return nil, domain.ErrPhoneConflict
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.q.CreateUser(ctx, repository.CreateUserParams{
		ID:           domain.NewUUIDV7(),
		Email:        pgtype.Text{String: req.Email, Valid: true},
		Phone:        req.Phone,
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
		FullName:     req.FullName,
		Role:         role,
		IsVerified:   pgtype.Bool{Bool: false, Valid: true},
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userIDStr := domain.UUIDToString(user.ID)

	b := make([]byte, 16)
	_, _ = rand.Read(b)
	verifyToken := hex.EncodeToString(b)

	s.rdb.Set(ctx, fmt.Sprintf("auth_verify_token:%s", verifyToken), userIDStr, 24*time.Hour)
	log.Printf("[BACKGROUND JOB] Enqueued verification email with token %s to %s", verifyToken, req.Email)

	_, _ = s.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		ID:     domain.NewUUIDV7(),
		UserID:    user.ID,
		Action:    "register",
		IpAddress: ip,
		UserAgent: userAgent,
	})

	return &domain.RegisterResp{
		UserID:   userIDStr,
		Email:    user.Email.String,
		Phone:    user.Phone,
		FullName: user.FullName,
		Role:     user.Role,
	}, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginReq, ip, userAgent string) (*domain.LoginResp, error) {
	blocked, ttl, err := s.limiter.IsIPBlocked(ctx, ip)
	if err != nil {
		log.Printf("[WARNING] Rate limiter error: %v", err)
	}
	if blocked {
		return nil, fmt.Errorf("%w: %d", domain.ErrIPBlocked, int(ttl.Seconds()))
	}

	if req.Email == "" || req.Password == "" {
		return nil, domain.NewValidationError("email and password are required")
	}

	user, err := s.q.GetUserByEmail(ctx, pgtype.Text{String: req.Email, Valid: true})
	if err != nil {
		_, _ = s.limiter.RecordFailedLogin(ctx, ip)
		return nil, domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password))
	if err != nil {
		_, _ = s.limiter.RecordFailedLogin(ctx, ip)
		return nil, domain.ErrInvalidCredentials
	}

	_ = s.limiter.ResetFailedLogin(ctx, ip)

	userIDStr := domain.UUIDToString(user.ID)

	accessToken, refreshToken, err := utils.GenerateTokenPair(userIDStr, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	err = s.rdb.Set(ctx, fmt.Sprintf("user:%s:refresh", userIDStr), refreshToken, 30*24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	_, _ = s.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		ID:     domain.NewUUIDV7(),
		UserID:    user.ID,
		Action:    "login",
		IpAddress: ip,
		UserAgent: userAgent,
	})

	return &domain.LoginResp{
		UserID:       userIDStr,
		FullName:     user.FullName,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) GoogleLogin(ctx context.Context, req domain.GoogleLoginReq, ip, userAgent string) (*domain.LoginResp, error) {
	gUser, err := utils.VerifyGoogleToken(ctx, req.Token)
	if err != nil {
		return nil, domain.NewValidationError(err.Error())
	}

	var user repository.User

	user, err = s.q.GetUserByGoogleID(ctx, pgtype.Text{String: gUser.ID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			existingUser, emailErr := s.q.GetUserByEmail(ctx, pgtype.Text{String: gUser.Email, Valid: true})
			if emailErr == nil {
				user, err = s.q.UpdateUserGoogleID(ctx, repository.UpdateUserGoogleIDParams{
					ID:       existingUser.ID,
					GoogleID: pgtype.Text{String: gUser.ID, Valid: true},
				})
				if err != nil {
					return nil, fmt.Errorf("failed to link Google account: %w", err)
				}
				_, _ = s.q.UpdateUserVerificationStatus(ctx, repository.UpdateUserVerificationStatusParams{
					ID:         user.ID,
					IsVerified: pgtype.Bool{Bool: true, Valid: true},
				})
			} else {
				uniqueMockPhone := "+62899" + gUser.ID[len(gUser.ID)-8:]
				user, err = s.q.CreateUser(ctx, repository.CreateUserParams{
					ID:         domain.NewUUIDV7(),
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
					return nil, fmt.Errorf("failed to register Google account: %w", err)
				}
			}
		} else {
			return nil, fmt.Errorf("failed to query user: %w", err)
		}
	}

	userIDStr := domain.UUIDToString(user.ID)

	accessToken, refreshToken, err := utils.GenerateTokenPair(userIDStr, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.rdb.Set(ctx, fmt.Sprintf("user:%s:refresh", userIDStr), refreshToken, 30*24*time.Hour)

	_, _ = s.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		ID:     domain.NewUUIDV7(),
		UserID:    user.ID,
		Action:    "login_google",
		IpAddress: ip,
		UserAgent: userAgent,
	})

	return &domain.LoginResp{
		UserID:       userIDStr,
		FullName:     user.FullName,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) SendWhatsAppOTP(ctx context.Context, req domain.SendOtpReq) error {
	if !phoneRegex().MatchString(req.Phone) {
		return domain.NewValidationError("invalid Indonesian phone number format")
	}

	otp := utils.GenerateOTP()

	err := s.rdb.Set(ctx, fmt.Sprintf("whatsapp_otp:%s", req.Phone), otp, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	err = utils.SendWhatsAppOTP(ctx, req.Phone, otp)
	if err != nil {
		return fmt.Errorf("failed to dispatch WhatsApp OTP: %w", err)
	}

	return nil
}

func (s *authService) VerifyWhatsAppOTP(ctx context.Context, req domain.VerifyOtpReq, ip, userAgent string) (*domain.LoginResp, error) {
	if req.Phone == "" || req.OTP == "" {
		return nil, domain.NewValidationError("phone number and OTP are required")
	}

	otpKey := fmt.Sprintf("whatsapp_otp:%s", req.Phone)

	storedOtp, err := s.rdb.Get(ctx, otpKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrIncorrectOTP
		}
		return nil, fmt.Errorf("failed to verify OTP: %w", err)
	}

	if storedOtp != req.OTP {
		return nil, domain.ErrIncorrectOTP
	}

	s.rdb.Del(ctx, otpKey)

	var user repository.User

	user, err = s.q.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			user, err = s.q.CreateUser(ctx, repository.CreateUserParams{
				ID:         domain.NewUUIDV7(),
				Phone:      req.Phone,
				FullName:   "WhatsApp User",
				Role:       "user",
				IsVerified: pgtype.Bool{Bool: true, Valid: true},
				IsActive:   pgtype.Bool{Bool: true, Valid: true},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to auto-register user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to query user by phone: %w", err)
		}
	}

	userIDStr := domain.UUIDToString(user.ID)

	accessToken, refreshToken, err := utils.GenerateTokenPair(userIDStr, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.rdb.Set(ctx, fmt.Sprintf("user:%s:refresh", userIDStr), refreshToken, 30*24*time.Hour)

	_, _ = s.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		ID:     domain.NewUUIDV7(),
		UserID:    user.ID,
		Action:    "login_whatsapp",
		IpAddress: ip,
		UserAgent: userAgent,
	})

	return &domain.LoginResp{
		UserID:       userIDStr,
		FullName:     user.FullName,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req domain.RefreshTokenReq) (*domain.RefreshTokenResp, error) {
	if req.RefreshToken == "" {
		return nil, domain.NewValidationError("refresh token is required")
	}

	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, domain.NewValidationError("invalid or expired refresh token")
	}

	redisKey := fmt.Sprintf("user:%s:refresh", claims.UserID)

	storedToken, err := s.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		return nil, domain.NewValidationError("session has expired or logout completed. Please login again.")
	}

	if storedToken != req.RefreshToken {
		s.rdb.Del(ctx, redisKey)
		return nil, domain.NewValidationError("security warning: refresh token mismatch. Session revoked.")
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(claims.UserID, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to rotate tokens: %w", err)
	}

	s.rdb.Set(ctx, redisKey, refreshToken, 30*24*time.Hour)

	return &domain.RefreshTokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	s.rdb.Del(ctx, fmt.Sprintf("user:%s:refresh", userID))
	return nil
}

func (s *authService) ForgotPassword(ctx context.Context, req domain.ForgotPasswordReq) error {
	if req.Email == "" {
		return domain.NewValidationError("email address is required")
	}

	user, err := s.q.GetUserByEmail(ctx, pgtype.Text{String: req.Email, Valid: true})
	if err != nil {
		return nil
	}

	b := make([]byte, 20)
	_, _ = rand.Read(b)
	resetToken := hex.EncodeToString(b)

	userIDStr := domain.UUIDToString(user.ID)
	s.rdb.Set(ctx, fmt.Sprintf("reset_password_token:%s", resetToken), userIDStr, 1*time.Hour)

	log.Printf("[BACKGROUND JOB] Enqueued password reset email with token %s to %s", resetToken, req.Email)

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req domain.ResetPasswordReq, ip, userAgent string) error {
	if req.Token == "" || req.NewPassword == "" {
		return domain.NewValidationError("reset token and new password are required")
	}

	if len(req.NewPassword) < 8 {
		return domain.NewValidationError("new password must be at least 8 characters long")
	}

	tokenKey := fmt.Sprintf("reset_password_token:%s", req.Token)

	userIDStr, err := s.rdb.Get(ctx, tokenKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.NewValidationError("reset token is invalid or has expired")
		}
		return fmt.Errorf("failed to verify reset token: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	userID := domain.StringToUUID(userIDStr)
	_, err = s.q.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.rdb.Del(ctx, fmt.Sprintf("user:%s:refresh", userIDStr))
	s.rdb.Del(ctx, tokenKey)

	_, _ = s.q.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		ID:     domain.NewUUIDV7(),
		UserID:    userID,
		Action:    "reset_password",
		IpAddress: ip,
		UserAgent: userAgent,
	})

	return nil
}
