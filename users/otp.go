package users

import (
	"github.com/pquerna/otp/totp"
)

func CheckOtp(otp, authenticatorToken string) bool {
	// OTP validation: only enforce if AUTHENTICATOR_TOKEN environment variable is set.
	if authenticatorToken == "" {
		return true
	}
	// Validate the TOTP code using the TOTP library
	return totp.Validate(otp, authenticatorToken)
}
