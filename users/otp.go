package users

import (
	"github.com/pquerna/otp/totp"
)

func CheckOtp(otp, authenticatorToken string) bool {
	// OTP validation: only enforce if authenticatorToken is set in the database (filebrowser config)
	if authenticatorToken == "" {
		return true
	}
	// Validate the TOTP code using the TOTP library
	return totp.Validate(otp, authenticatorToken)
}
