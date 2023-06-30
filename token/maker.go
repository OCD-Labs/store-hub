package token

import "time"

// A Maker is an interface for managing tokens.
type Maker interface {
	// CreateToken creates a new token for a specific username and duration.
	CreateToken(userID int64, account_id, userRole string, permissionLevel *int16, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}