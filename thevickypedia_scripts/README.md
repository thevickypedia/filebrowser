## Primary Changes

1. JSON authentication

`frontend/src/utils/auth.ts` - Set authorization header
`auth/json.go` - Decrypt authorization header
`auth/database.go` - Handle auth errors
`http/connection.go` - Handles connection logging

2. Handle allowed origins

`cmd/root.go` - Background task
`cmd/ip_addresses.go` - Refresh allowed origins in the background

3. OTP

`users/otp.go` - Verifies the one-time passcode.
