// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Item struct {
	ID                 int64           `json:"id"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Price              string          `json:"price"`
	StoreID            int64           `json:"store_id"`
	ImageUrls          []string        `json:"image_urls"`
	Category           string          `json:"category"`
	DiscountPercentage string          `json:"discount_percentage"`
	SupplyQuantity     int64           `json:"supply_quantity"`
	Extra              json.RawMessage `json:"extra"`
	IsFrozen           bool            `json:"is_frozen"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	Currency           string          `json:"currency"`
	CoverImgUrl        string          `json:"cover_img_url"`
}

type ItemRating struct {
	UserID    int64          `json:"user_id"`
	ItemID    int64          `json:"item_id"`
	Rating    int16          `json:"rating"`
	Comment   sql.NullString `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
}

type Order struct {
	ID                   int64     `json:"id"`
	DeliveryStatus       string    `json:"delivery_status"`
	DeliveredOn          time.Time `json:"delivered_on"`
	ExpectedDeliveryDate time.Time `json:"expected_delivery_date"`
	ItemID               int64     `json:"item_id"`
	OrderQuantity        int32     `json:"order_quantity"`
	BuyerID              int64     `json:"buyer_id"`
	SellerID             int64     `json:"seller_id"`
	StoreID              int64     `json:"store_id"`
	DeliveryFee          string    `json:"delivery_fee"`
	PaymentChannel       string    `json:"payment_channel"`
	PaymentMethod        string    `json:"payment_method"`
	CreatedAt            time.Time `json:"created_at"`
}

type Sale struct {
	ID         int64     `json:"id"`
	StoreID    int64     `json:"store_id"`
	ItemID     int64     `json:"item_id"`
	CustomerID int64     `json:"customer_id"`
	SellerID   int64     `json:"seller_id"`
	OrderID    int64     `json:"order_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type SalesOverview struct {
	ID              int64  `json:"id"`
	NumberOfSales   int64  `json:"number_of_sales"`
	SalesPercentage string `json:"sales_percentage"`
	Revenue         string `json:"revenue"`
	ItemID          int64  `json:"item_id"`
	StoreID         int64  `json:"store_id"`
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
	StoreAccountID  string    `json:"store_account_id"`
	ProfileImageUrl string    `json:"profile_image_url"`
	IsVerified      bool      `json:"is_verified"`
	Category        string    `json:"category"`
	IsFrozen        bool      `json:"is_frozen"`
	CreatedAt       time.Time `json:"created_at"`
}

type StoreAuditTrail struct {
	ID        int64                 `json:"id"`
	StoreID   int64                 `json:"store_id"`
	UserID    int64                 `json:"user_id"`
	Action    string                `json:"action"`
	Details   pqtype.NullRawMessage `json:"details"`
	Timestamp time.Time             `json:"timestamp"`
}

type StoreOwner struct {
	UserID       int64     `json:"user_id"`
	StoreID      int64     `json:"store_id"`
	AddedAt      time.Time `json:"added_at"`
	IsPrimary    bool      `json:"is_primary"`
	AccessLevels []int32   `json:"access_levels"`
}

type User struct {
	ID                int64           `json:"id"`
	FirstName         string          `json:"first_name"`
	LastName          string          `json:"last_name"`
	AccountID         string          `json:"account_id"`
	Status            string          `json:"status"`
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
