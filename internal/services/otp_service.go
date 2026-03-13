package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"time"

	"hisabi.com/m/internal/repository"
)

func generateOtp() (string, error) {
	b := make([]byte, 3)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	otp := fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2]))
	if len(otp) > 6 {
		otp = otp[len(otp)-6:]
	}
	return otp, nil
}

func SendVerificationOTP(email, username string) error {

	otp, err := generateOtp()
	if err != nil {
		return errors.New("failed to generate OTP")
	}

	expiresAt := time.Now().Add(10 * time.Minute)

	if err := repository.SaveOTP(email, otp, expiresAt); err != nil {
		return errors.New("failed to save OTP")
	}

	if err := SendOTPEmail(email, username, otp); err != nil {
		log.Printf("Email send failed | Email: %s | Err: %v", email, err)
		return errors.New("failed to send verification email. Please try again")
	}

	log.Printf("OTP sent | Email: %s | Expires: %s", email, expiresAt.Format("15:04:05"))
	return nil
}

func VerifyOTP(email, code string) error {

	// Input check
	if email == "" || code == "" {
		return errors.New("email and OTP code are required")
	}
	if len(code) != 6 {
		return errors.New("OTP must be 6 digits")
	}

	otp, err := repository.GetValidOTP(email, code)
	if err != nil {

		return errors.New("invalid or expired OTP")
	}

	if err := repository.MarkOTPUsed(otp.ID); err != nil {
		return errors.New("verification failed, please try again")
	}

	if err := repository.MarkUserVerified(email); err != nil {
		return errors.New("failed to verify account")
	}

	log.Printf("Email verified | Email: %s", email)
	return nil
}

func ResendOTP(email string) error {

	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return errors.New("no account found with this email")
	}

	if user.IsVerified {
		return errors.New("this account is already verified")
	}
	return SendVerificationOTP(email, user.Username)
}
