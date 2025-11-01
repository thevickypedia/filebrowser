package users

import "os"

func CheckOtp(otp string) bool {
	// OTP validation: only enforce if AUTHENTICATOR_TOKEN environment variable is set.
	//nolint:godox // reason: placeholder until real TOTP added
	// TODO: This should be using TOTP library for real OTP validation.
	envOtp := os.Getenv("AUTHENTICATOR_TOKEN")
	if envOtp != "" {
		return otp == envOtp
	}
	return true
}
