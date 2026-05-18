package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go-beresin/pkg/utils"
)

// Response represents a standard JSON response structure.
type Response struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginRequest defines the request body for mock login.
type LoginRequest struct {
	Username string `json:"username" example:"johndoe" xml:"username" form:"username"`
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
// @Success 200 {object} Response
// @Router /health [get]
func (h *DummyHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Server is healthy and running",
	})
}

// LoginMock handles simulated authentication for test purposes.
// @Summary Mock Login
// @Description Authenticate a dummy user and return a JWT token (roles: user, admin)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Payload"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /api/v1/auth/login [post]
func (h *DummyHandler) LoginMock(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	if req.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Username is required",
		})
	}

	// Default role to user if empty
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Generate standard JWT token
	token, err := utils.GenerateToken(req.Username, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: fmt.Sprintf("Failed to generate token: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully authenticated",
		Data: fiber.Map{
			"username": req.Username,
			"role":     role,
			"token":    token,
		},
	})
}

// Profile handles retrieving authenticated user details from local context.
// @Summary Get Profile
// @Description Fetch authenticated user identity from Bearer token claims
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /api/v1/profile [get]
func (h *DummyHandler) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	role := c.Locals("role")

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Profile retrieved successfully",
		Data: fiber.Map{
			"user_id": userID,
			"role":    role,
		},
	})
}

// AdminOnly handles endpoints restricted to administrators.
// @Summary Admin Endpoint
// @Description Access administrator dashboard (restricted to role: admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Router /api/v1/admin [get]
func (h *DummyHandler) AdminOnly(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Welcome to the Admin Dashboard",
		Data: fiber.Map{
			"admin_user_id": userID,
			"privileges":    "unlimited",
		},
	})
}

// Upload handles receiving and saving a multipart file.
// @Summary File Upload
// @Description Upload any photo, video, or document (max 10MB limit)
// @Tags Uploads
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/upload [post]
func (h *DummyHandler) Upload(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "No file was uploaded or file parameter 'file' is missing",
		})
	}

	// Save to uploads folder
	savedPath, err := utils.ProcessUploadedFile(fileHeader, "uploads")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:  "success",
		Message: "File uploaded and verified successfully",
		Data: fiber.Map{
			"file_name": fileHeader.Filename,
			"file_size": fileHeader.Size,
			"save_path": savedPath,
		},
	})
}
