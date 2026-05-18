package domain

import "errors"

type RegisterReq struct {
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,phone"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"omitempty,oneof=user talent"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GoogleLoginReq struct {
	Token string `json:"token" validate:"required"`
}

type SendOtpReq struct {
	Phone string `json:"phone" validate:"required,phone"`
}

type VerifyOtpReq struct {
	Phone string `json:"phone" validate:"required,phone"`
	OTP   string `json:"otp" validate:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordReq struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type RegisterResp struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type LoginResp struct {
	UserID       string `json:"user_id"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(msg string) error {
	return &ValidationError{Message: msg}
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailConflict      = errors.New("email already registered")
	ErrPhoneConflict      = errors.New("phone already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrIPBlocked          = errors.New("ip temporarily blocked due to multiple failed login attempts")
	ErrIncorrectOTP       = errors.New("incorrect or expired OTP")
)
