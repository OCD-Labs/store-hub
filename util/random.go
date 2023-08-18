package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

const (
	STOREOWNER = "STORE-OWNER"
	NORMALUSER = "NORMAL-USER"
)

var DELIVERYSTATUS = []string{"PENDING", "PROCESSING", "SHIPPED", "DELIVERED", "CANCELLED", "RETURNED"}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates random integer between min and max.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// RandomPermission generates a random permission code.
func RandomPermission() string {
	permissions := []string{STOREOWNER, NORMALUSER}
	n := len(permissions)
	return permissions[rand.Intn(n)]
}

// RandomEmail generates a random email.
func RandomEmail() string {
	return fmt.Sprintf("%s-%s@gmail.com", RandomString(6), RandomString(6))
}

// Extract retrieve a substring of the PASETO token string value.
func Extract(s string) string {
	start := "v2.local."
	end := ".bnVsbA"
	startIndex := strings.Index(s, start)
	endIndex := strings.Index(s, end)

	if startIndex == -1 || endIndex == -1 {
		return ""
	}

	startIndex += len(start)
	return s[startIndex:endIndex]
}

// Concat concatenates the substring of the PASETO token string value.
func Concat(s string) string {
	return fmt.Sprintf("v2.local.%s.bnVsbA", s)
}

// IsValidStatus check if status exists in DELIVERYSTATUS
func IsValidStatus(status string) bool {
	status = strings.TrimSpace(status)
	for _, s := range DELIVERYSTATUS {
		if s == status {
			return true
		}
	}
	return false
}

// ConvertToPercentage converts a float64 value to a formatted percentage string.
func ConvertToPercentage(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

// CanChangeStatus checks if the nextStatus is one of the allowed statuses for that currentStatus.
func CanChangeStatus(currentStatus, nextStatus string) bool {
	switch currentStatus {
	case "PENDING":
		return nextStatus == "PROCESSING" || nextStatus == "SHIPPED" || nextStatus == "DELIVERED" || nextStatus == "CANCELLED"
	case "PROCESSING":
		return nextStatus == "SHIPPED" || nextStatus == "DELIVERED" || nextStatus == "CANCELLED"
	case "SHIPPED":
		return nextStatus == "DELIVERED" || nextStatus == "CANCELLED"
	case "DELIVERED":
		return nextStatus == "CANCELLED"
	case "CANCELLED":
		return nextStatus == "RETURNED"
	default:
		return false
	}
}
