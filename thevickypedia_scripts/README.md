## Primary Changes

1. JSON authentication

`frontend/src/utils/auth.ts` - Set authorization header
`auth/json.go` - Decrypt authorization header
`auth/database.go` - Handle auth errors

2. Handle allowed origins

`cmd/root.go` - Background task
`cmd/ip_addresses.go` - Refresh allowed origins in the background
