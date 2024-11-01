// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Cart struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CartItem struct {
	ID        int64     `json:"id"`
	CartID    int64     `json:"cart_id"`
	ItemID    int64     `json:"item_id"`
	StoreID   int64     `json:"store_id"`
	Quantity  int32     `json:"quantity"`
	AddedAt   time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CryptoAccount struct {
	ID            int64     `json:"id"`
	StoreID       int64     `json:"store_id"`
	Balance       string    `json:"balance"`
	WalletAddress string    `json:"wallet_address"`
	CryptoType    string    `json:"crypto_type"`
	CreatedAt     time.Time `json:"created_at"`
}

type FiatAccount struct {
	ID        int64     `json:"id"`
	StoreID   int64     `json:"store_id"`
	Balance   string    `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

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
	Currency           string          `json:"currency"`
	CoverImgUrl        string          `json:"cover_img_url"`
	Status             string          `json:"status"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type Order struct {
	ID                   int64     `json:"id"`
	DeliveryStatus       string    `json:"delivery_status"`
	DeliveredOn          time.Time `json:"delivered_on"`
	ExpectedDeliveryDate time.Time `json:"expected_delivery_date"`
	ItemID               int64     `json:"item_id"`
	ItemPrice            string    `json:"item_price"`
	ItemCurrency         string    `json:"item_currency"`
	OrderQuantity        int32     `json:"order_quantity"`
	BuyerID              int64     `json:"buyer_id"`
	SellerID             int64     `json:"seller_id"`
	StoreID              int64     `json:"store_id"`
	DeliveryFee          string    `json:"delivery_fee"`
	PaymentChannel       string    `json:"payment_channel"`
	PaymentMethod        string    `json:"payment_method"`
	IsReviewed           bool      `json:"is_reviewed"`
	CreatedAt            time.Time `json:"created_at"`
}

type PendingTransactionFund struct {
	ID          int64     `json:"id"`
	StoreID     int64     `json:"store_id"`
	AccountType string    `json:"account_type"`
	Amount      string    `json:"amount"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type Review struct {
	ID                 int64     `json:"id"`
	StoreID            int64     `json:"store_id"`
	UserID             int64     `json:"user_id"`
	ItemID             int64     `json:"item_id"`
	Rating             string    `json:"rating"`
	ReviewType         string    `json:"review_type"`
	Comment            string    `json:"comment"`
	IsVerifiedPurchase bool      `json:"is_verified_purchase"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ReviewLike struct {
	ID       int64 `json:"id"`
	ReviewID int64 `json:"review_id"`
	UserID   int64 `json:"user_id"`
	Liked    bool  `json:"liked"`
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
	AccessLevels []int32   `json:"access_levels"`
	IsPrimary    bool      `json:"is_primary"`
	AddedAt      time.Time `json:"added_at"`
}

type Transaction struct {
	ID                   int64          `json:"id"`
	OrderIds             []int64        `json:"order_ids"`
	CustomerID           int64          `json:"customer_id"`
	Amount               string         `json:"amount"`
	PaymentProvider      string         `json:"payment_provider"`
	ProviderTxRefID      string         `json:"provider_tx_ref_id"`
	ProviderTxAccessCode sql.NullString `json:"provider_tx_access_code"`
	ProviderTxFee        string         `json:"provider_tx_fee"`
	Status               string         `json:"status"`
	CreatedAt            time.Time      `json:"created_at"`
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
