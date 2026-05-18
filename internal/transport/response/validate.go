package response

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

var phoneRegex = regexp.MustCompile(`^(?:\+62|62|0)8[1-9][0-9]{6,10}$`)

func init() {
	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		return phoneRegex.MatchString(fl.Field().String())
	})
}

type validationErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func msgForTag(tag string, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least 8 characters", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of 'user' or 'talent'", field)
	case "phone":
		return "invalid Indonesian phone number format"
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func fieldName(structField string) string {
	switch structField {
	case "Email":
		return "Email"
	case "Phone":
		return "Phone"
	case "Password":
		return "Password"
	case "FullName":
		return "Full Name"
	case "Role":
		return "Role"
	case "Token":
		return "Token"
	case "OTP":
		return "OTP"
	case "RefreshToken":
		return "Refresh Token"
	case "NewPassword":
		return "New Password"
	default:
		return structField
	}
}

func Check(c *fiber.Ctx, s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return ValidationError(c, "Validation failed", nil)
	}

	errs := make([]validationErr, 0, len(ve))
	for _, fe := range ve {
		errs = append(errs, validationErr{
			Field:   fe.Field(),
			Message: msgForTag(fe.Tag(), fieldName(fe.Field())),
		})
	}

	return ValidationError(c, "Validation failed", errs)
}
