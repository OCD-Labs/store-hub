// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID                 int64           `json:"id"`
	Description        string          `json:"description"`
	Price              interface{}     `json:"price"`
	StoreID            int64           `json:"store_id"`
	ImageUrls          []string        `json:"image_urls"`
	Category           string          `json:"category"`
	DiscountPercentage interface{}     `json:"discount_percentage"`
	SupplyQuantity     int64           `json:"supply_quantity"`
	Extra              json.RawMessage `json:"extra"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type ItemRating struct {
	UserID    int64          `json:"user_id"`
	ItemID    int64          `json:"item_id"`
	Rating    string         `json:"rating"`
	Comment   sql.NullString `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	Scope     string    `json:"scope"`
	UserAgent string    `json:"user_agent"`
	ClientIp  string    `json:"client_ip"`
	IsBlocked bool      `json:"is_blocked"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ProfileImageUrl string    `json:"profile_image_url"`
	IsVerified      bool      `json:"is_verified"`
	Category        string    `json:"category"`
	CreatedAt       time.Time `json:"created_at"`
}

type StoreOwner struct {
	UserID          int64     `json:"user_id"`
	StoreID         int64     `json:"store_id"`
	PermissionLevel int16     `json:"permission_level"`
	AddedAt         time.Time `json:"added_at"`
}

type User struct {
	ID                int64           `json:"id"`
	FirstName         string          `json:"first_name"`
	LastName          string          `json:"last_name"`
	Permission        string          `json:"permission"`
	About             string          `json:"about"`
	Email             string          `json:"email"`
	Socials           json.RawMessage `json:"socials"`
	ProfileImageUrl   sql.NullString  `json:"profile_image_url"`
	HashedPassword    string          `json:"hashed_password"`
	PasswordChangedAt time.Time       `json:"password_changed_at"`
	CreatedAt         time.Time       `json:"created_at"`
	IsActive          bool            `json:"is_active"`
	IsEmailVerified   bool            `json:"is_email_verified"`
}
