package util

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	_ = iota // ignore the first value
	FULLACCESS
	PRODUCTINVENTORYACCESS
	SALESACCESS
	ORDERSACCESS
	FINANCIALACCESS
	alphabets  = "abcdefghijklmnopqrstuvwxyz"
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

// NumberExists checks if number exist in the slice
func NumberExists(slice []int32, number int) bool {
	for _, v := range slice {
		if v == int32(number) {
			return true
		}
	}
	return false
}

// CommandExists checks if an executable named file exists
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// FolderExists checks if folder/file exists.
func FolderExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// SanitizeAccountID removes the suffix ".near" or ".testnet" from accountID.
func SanitizeAccountID(accountID, network string) string {
	switch network {
	case "near":
		accountID = strings.TrimSuffix(accountID, ".near")
	case "testnet":
		accountID = strings.TrimSuffix(accountID, ".testnet")
	}
	accountID = strings.TrimSuffix(accountID, ".storehub-v1")

	if len(accountID) == 0 {
		return RandomString(12)
	}

	// Replace invalid characters with underscores
	reg := regexp.MustCompile(`[^A-Za-z0-9_-]+`)
	return reg.ReplaceAllString(accountID, "_")
}


// EncodeToBase64 takes a byte slice as input and returns its Base64 encoded string
func EncodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
