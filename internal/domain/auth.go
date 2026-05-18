package domain

import "errors"

type RegisterReq struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
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
