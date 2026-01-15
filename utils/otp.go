package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type OTPStore struct {
	Code      string
	ExpiresAt time.Time
}

var otpStorage = make(map[string]OTPStore)

func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	
	otp := fmt.Sprintf("%06d", n.Int64())
	return otp, nil
}

func StoreOTP(bookingID string, otpType string) (string, error) {
	otp, err := GenerateOTP()
	if err != nil {
		return "", err
	}
	
	key := fmt.Sprintf("%s_%s", bookingID, otpType)
	otpStorage[key] = OTPStore{
		Code:      otp,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	
	return otp, nil
}

func VerifyOTP(bookingID string, otpType string, providedOTP string) bool {
	key := fmt.Sprintf("%s_%s", bookingID, otpType)
	
	stored, exists := otpStorage[key]
	if !exists {
		return false
	}
	
	if time.Now().After(stored.ExpiresAt) {
		delete(otpStorage, key)
		return false
	}
	
	if stored.Code != providedOTP {
		return false
	}
	
	delete(otpStorage, key)
	return true
}

func ClearOTP(bookingID string, otpType string) {
	key := fmt.Sprintf("%s_%s", bookingID, otpType)
	delete(otpStorage, key)
}
