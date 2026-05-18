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

	"go-beresin/pkg/database"
)

// Helper: Setup test dependencies
func setupTestDeps(t *testing.T) (*redis.Client, *AuthHandler, func()) {
	_ = godotenv.Load("../../../.env") // Load root .env if present

	ctx := context.Background()

	// Initialize DB pool
	dbPool, err := database.InitPool(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize database pool for testing: %v", err)
	}

	// Initialize Redis client
	redisAddr := "localhost:6379"
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "!Abcd1234",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to Redis for testing: %v", err)
	}

	// Ensure clean state before running tests (clear Fiber test client IP: 0.0.0.0)
	rdb.Del(ctx, "auth_failed_counter:127.0.0.1")
	rdb.Del(ctx, "auth_block:127.0.0.1")
	rdb.Del(ctx, "auth_failed_counter:0.0.0.0")
	rdb.Del(ctx, "auth_block:0.0.0.0")

	authH := NewAuthHandler(dbPool, rdb)

	cleanup := func() {
		rdb.Del(ctx, "auth_failed_counter:127.0.0.1")
		rdb.Del(ctx, "auth_block:127.0.0.1")
		rdb.Del(ctx, "auth_failed_counter:0.0.0.0")
		rdb.Del(ctx, "auth_block:0.0.0.0")
		dbPool.Close()
		rdb.Close()
	}

	return rdb, authH, cleanup
}

func TestAuthFlow(t *testing.T) {
	rdb, authH, cleanup := setupTestDeps(t)
	defer cleanup()

	// 1. Setup test Fiber App
	app := fiber.New()
	app.Post("/api/v1/auth/register", authH.Register)
	app.Post("/api/v1/auth/login", authH.Login)
	app.Post("/api/v1/auth/refresh-token", authH.RefreshToken)

	// Generate random email and phone to guarantee unique tests
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(1000000)
	testEmail := fmt.Sprintf("test_auth_%d@example.com", randomInt)
	testPhone := fmt.Sprintf("+62812345%05d", rand.Intn(100000))
	testPassword := "supersecret123"

	// 2. TEST REGISTER - SUCCESS
	registerPayload := RegisterReq{
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

	// 3. TEST REGISTER - CONFLICT (DUPLICATE EMAIL)
	reqConflict := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	reqConflict.Header.Set("Content-Type", "application/json")
	respConflict, err := app.Test(reqConflict, -1)
	if err != nil {
		t.Fatalf("Failed to dispatch duplicate Register request: %v", err)
	}
	if respConflict.StatusCode != http.StatusConflict {
		t.Errorf("Expected registration conflict status 409, got %d", respConflict.StatusCode)
	}

	// 4. TEST LOGIN - SUCCESS
	loginPayload := LoginReq{
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

	// 5. TEST REFRESH TOKEN - SUCCESS
	refreshPayload := RefreshTokenReq{
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

	// 6. Clean up test user in DB after completing assertions
	ctx := context.Background()
	_, _ = authH.db.Exec(ctx, "DELETE FROM users WHERE email = $1", testEmail)

	// Clean up Redis refresh token session key
	_, _ = rdb.Del(ctx, fmt.Sprintf("user:%s:refresh", loginResponse.Data.AccessToken)).Result()
}

func TestLoginRateLimiting(t *testing.T) {
	rdb, authH, cleanup := setupTestDeps(t)
	defer cleanup()

	// Clear failed login counter key in Redis to start clean
	ctx := context.Background()
	mockIP := "0.0.0.0" // Fiber default test client IP
	rdb.Del(ctx, fmt.Sprintf("auth_failed_counter:%s", mockIP))
	rdb.Del(ctx, fmt.Sprintf("auth_block:%s", mockIP))

	app := fiber.New()
	app.Post("/api/v1/auth/login", authH.Login)

	loginPayload := LoginReq{
		Email:    "nonexistent_user@example.com",
		Password: "wrongpassword",
	}
	loginBody, _ := json.Marshal(loginPayload)

	// Fail login 5 times to trigger temporary IP block
	for i := 1; i <= 6; i++ {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("Failed on iteration %d: %v", i, err)
		}

		if i <= 5 {
			// Should return 401 Unauthorized
			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("Iteration %d: Expected status 401, got %d", i, resp.StatusCode)
			}
		} else {
			// 6th and subsequent attempts should return 429 Too Many Requests
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("Iteration %d: Expected status 429, got %d", i, resp.StatusCode)
			}
		}
	}

	// Clean up Redis
	rdb.Del(ctx, fmt.Sprintf("auth_failed_counter:%s", mockIP))
	rdb.Del(ctx, fmt.Sprintf("auth_block:%s", mockIP))
}
