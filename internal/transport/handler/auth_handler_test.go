package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"go-beresin/internal/domain"
	"go-beresin/internal/service"
	"go-beresin/pkg/database"
)

func setupTestDeps(t *testing.T) (*redis.Client, *AuthHandler, *service.AuthService, func()) {
	_ = godotenv.Load("../../../.env")

	ctx := context.Background()

	dbPool, err := database.InitPool(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize database pool for testing: %v", err)
	}

	redisAddr := "localhost:6379"
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "!Abcd1234",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to Redis for testing: %v", err)
	}

	rdb.Del(ctx, "auth_failed_counter:127.0.0.1")
	rdb.Del(ctx, "auth_block:127.0.0.1")
	rdb.Del(ctx, "auth_failed_counter:0.0.0.0")
	rdb.Del(ctx, "auth_block:0.0.0.0")

	authSvc := service.NewAuthService(dbPool, rdb)
	authH := NewAuthHandler(authSvc)

	cleanup := func() {
		rdb.Del(ctx, "auth_failed_counter:127.0.0.1")
		rdb.Del(ctx, "auth_block:127.0.0.1")
		rdb.Del(ctx, "auth_failed_counter:0.0.0.0")
		rdb.Del(ctx, "auth_block:0.0.0.0")
		dbPool.Close()
		rdb.Close()
	}

	return rdb, authH, &authSvc, cleanup
}

func TestAuthFlow(t *testing.T) {
	rdb, authH, _, cleanup := setupTestDeps(t)
	defer cleanup()

	app := fiber.New()
	app.Post("/api/v1/auth/register", authH.Register)
	app.Post("/api/v1/auth/login", authH.Login)
	app.Post("/api/v1/auth/refresh-token", authH.RefreshToken)

	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(1000000)
	testEmail := fmt.Sprintf("test_auth_%d@example.com", randomInt)
	testPhone := fmt.Sprintf("+62812345%05d", rand.Intn(100000))
	testPassword := "supersecret123"

	registerPayload := domain.RegisterReq{
		Email:    testEmail,
		Phone:    testPhone,
		Password: testPassword,
		FullName: "Test Authentication Suite",
		Role:     "user",
	}
	body, _ := json.Marshal(registerPayload)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("Failed to dispatch Register request: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected registration status 201, got %d", resp.StatusCode)
	}

	reqConflict := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	reqConflict.Header.Set("Content-Type", "application/json")
	respConflict, err := app.Test(reqConflict, -1)
	if err != nil {
		t.Fatalf("Failed to dispatch duplicate Register request: %v", err)
	}
	if respConflict.StatusCode != http.StatusConflict {
		t.Errorf("Expected registration conflict status 409, got %d", respConflict.StatusCode)
	}

	loginPayload := domain.LoginReq{
		Email:    testEmail,
		Password: testPassword,
	}
	loginBody, _ := json.Marshal(loginPayload)
	reqLogin := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")

	respLogin, err := app.Test(reqLogin, -1)
	if err != nil {
		t.Fatalf("Failed to dispatch Login request: %v", err)
	}
	if respLogin.StatusCode != http.StatusOK {
		t.Errorf("Expected login status 200, got %d", respLogin.StatusCode)
	}

	var loginResponse struct {
		Status string `json:"status"`
		Data   struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(respLogin.Body).Decode(&loginResponse)

	if loginResponse.Data.AccessToken == "" || loginResponse.Data.RefreshToken == "" {
		t.Error("Expected valid AccessToken and RefreshToken in login response")
	}

	refreshPayload := domain.RefreshTokenReq{
		RefreshToken: loginResponse.Data.RefreshToken,
	}
	refreshBody, _ := json.Marshal(refreshPayload)
	reqRefresh := httptest.NewRequest("POST", "/api/v1/auth/refresh-token", bytes.NewBuffer(refreshBody))
	reqRefresh.Header.Set("Content-Type", "application/json")

	respRefresh, err := app.Test(reqRefresh, -1)
	if err != nil {
		t.Fatalf("Failed to dispatch Refresh request: %v", err)
	}
	if respRefresh.StatusCode != http.StatusOK {
		t.Errorf("Expected refresh status 200, got %d", respRefresh.StatusCode)
	}

	ctx := context.Background()
	_, _ = rdb.Del(ctx, fmt.Sprintf("user:%s:refresh", loginResponse.Data.AccessToken)).Result()
}

func TestLoginRateLimiting(t *testing.T) {
	_, authH, _, cleanup := setupTestDeps(t)
	defer cleanup()

	ctx := context.Background()
	mockIP := "0.0.0.0"

	app := fiber.New()
	app.Post("/api/v1/auth/login", authH.Login)

	loginPayload := domain.LoginReq{
		Email:    "nonexistent_user@example.com",
		Password: "wrongpassword",
	}
	loginBody, _ := json.Marshal(loginPayload)

	for i := 1; i <= 6; i++ {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("Failed on iteration %d: %v", i, err)
		}

		if i <= 5 {
			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("Iteration %d: Expected status 401, got %d", i, resp.StatusCode)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("Iteration %d: Expected status 429, got %d", i, resp.StatusCode)
			}
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "!Abcd1234",
		DB:       0,
	})
	rdb.Del(ctx, fmt.Sprintf("auth_failed_counter:%s", mockIP))
	rdb.Del(ctx, fmt.Sprintf("auth_block:%s", mockIP))
	rdb.Close()
}
