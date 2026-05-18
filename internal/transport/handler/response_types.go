package handler

import "go-beresin/internal/domain"

type RegisterResp struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Registration completed successfully. Please check your email to verify your account."`
	Data    domain.RegisterResp `json:"data"`
}

type LoginResp struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"Successfully authenticated"`
	Data    domain.LoginResp `json:"data"`
}

type RefreshTokenResp struct {
	Success bool                  `json:"success" example:"true"`
	Message string                `json:"message" example:"Tokens successfully rotated"`
	Data    domain.RefreshTokenResp `json:"data"`
}

type MessageResp struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
}

type ErrorResp struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Invalid request payload"`
}
