package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go-beresin/internal/transport/response"
	"go-beresin/pkg/utils"
)

// LoginRequest defines the request body for mock login.
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"johndoe" xml:"username" form:"username"`
	Role     string `json:"role" example:"admin" xml:"role" form:"role"`
}

// DummyHandler holds dependencies for handlers.
type DummyHandler struct{}

// NewDummyHandler creates a new instance of DummyHandler.
func NewDummyHandler() *DummyHandler {
	return &DummyHandler{}
}

// HealthCheck handles basic server health monitoring.
// @Summary Health Check
// @Description Check if the server is up and healthy
// @Tags System
// @Produce json
// @Success 200 {object} MessageResp
// @Router /health [get]
func (h *DummyHandler) HealthCheck(c *fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, "Server is healthy and running", nil)
}

// LoginMock handles simulated authentication for test purposes.
// @Summary Mock Login
// @Description Authenticate a dummy user and return a JWT token (roles: user, admin)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Payload"
// @Success 200 {object} MessageResp
// @Failure 400 {object} ErrorResp
// @Router /api/v1/auth/login-mock [post]
func (h *DummyHandler) LoginMock(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := response.Check(c, req); err != nil {
		return err
	}

	// Default role to user if empty
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Generate standard JWT token
	token, err := utils.GenerateToken(req.Username, role)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to generate token: %v", err))
	}

	return response.Success(c, fiber.StatusOK, "Successfully authenticated", fiber.Map{
		"username": req.Username,
		"role":     role,
		"token":    token,
	})
}

// AdminOnly handles endpoints restricted to administrators.
// @Summary Admin Endpoint
// @Description Access administrator dashboard (restricted to role: admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MessageResp
// @Failure 401 {object} ErrorResp
// @Failure 403 {object} ErrorResp
// @Router /api/v1/admin [get]
func (h *DummyHandler) AdminOnly(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	return response.Success(c, fiber.StatusOK, "Welcome to the Admin Dashboard", fiber.Map{
		"admin_user_id": userID,
		"privileges":    "unlimited",
	})
}

// Upload handles receiving and saving a multipart file.
// @Summary File Upload
// @Description Upload any photo, video, or document (max 10MB limit)
// @Tags Uploads
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} MessageResp
// @Failure 400 {object} ErrorResp
// @Failure 500 {object} ErrorResp
// @Router /api/v1/upload [post]
func (h *DummyHandler) Upload(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "No file was uploaded or file parameter 'file' is missing")
	}

	// Save to uploads folder
	savedPath, err := utils.ProcessUploadedFile(fileHeader, "uploads")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "File uploaded and verified successfully", fiber.Map{
		"file_name": fileHeader.Filename,
		"file_size": fileHeader.Size,
		"save_path": savedPath,
	})
}
