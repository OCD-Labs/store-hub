// Package token (payload) defines the token's payload.
package token

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidToken is returned for invalid token
	ErrInvalidToken = errors.New("token is invalid")

	// ErrExpiredToken is returned for expired token
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token.
type Payload struct {
	ID        uuid.UUID   `json:"id"`
	UserID    int64       `json:"user_id"`
	AccountID string      `json:"account_id"`
	Extra     interface{} `json:"extra"`
	IssuedAt  time.Time   `json:"issued_at"`
	ExpiredAt time.Time   `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration.
func NewPayload(userID int64, account_id string, duration time.Duration, extra interface{}) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		AccountID: account_id,
		Extra:     extra,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks for the expiry-validity of the token.
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}

// ExtractExtra extracts the Extra field from a given Payload into the provided interface.
func ExtractExtra(payload *Payload, target interface{}) error {
	// Step 1: Marshal payload.Extra back into a byte slice
	extraBytes, err := json.Marshal(payload.Extra)
	if err != nil {
		return errors.New("failed to marshal payload.Extra")
	}

	// Step 2: Unmarshal the byte slice into the provided target interface
	if err := json.Unmarshal(extraBytes, target); err != nil {
		return errors.New("failed to unmarshal into target interface")
	}

	return nil
}
