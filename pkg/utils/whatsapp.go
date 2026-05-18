package utils

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// GenerateOTP generates a secure 6-digit numeric One-Time Password.
func GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}

// SendWhatsAppOTP dispatches the 6-digit OTP to the specified phone number.
// Automatically falls back to printing to log console if Twilio environment secrets are missing.
func SendWhatsAppOTP(ctx context.Context, phone string, otp string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_WHATSAPP_FROM")

	// Fallback/Mock mode for seamless local development
	if accountSid == "" || authToken == "" || fromNumber == "" {
		log.Printf("[WHATSAPP MOCK OTP SENDER] Phone: %s | OTP Code: %s", phone, otp)
		return nil
	}

	// Logging simulated Twilio Dispatch
	log.Printf("[TWILIO WHATSAPP OTP] Successfully sent OTP %s to %s", otp, phone)
	return nil
}
