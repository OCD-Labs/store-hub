package accountid

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type AccountID string

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *AccountID) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Validate the account ID according to NEAR rules
	if !isValidAccountID(raw) {
		return fmt.Errorf("invalid account ID: %s", raw)
	}

	*a = AccountID(raw)
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (a AccountID) MarshalJSON() ([]byte, error) {
	// Check if the account ID is valid before marshaling
	if !isValidAccountID(string(a)) {
		return nil, fmt.Errorf("invalid account ID: %s", a)
	}

	return json.Marshal(string(a))
}

// isValidAccountID checks if an account ID follows NEAR rules
func isValidAccountID(accountID string) bool {
	// Check minimum and maximum length
	if len(accountID) < 2 || len(accountID) > 64 {
		return false
	}

	// Check that parts are alphanumeric and separated by _ or -
	parts := strings.Split(accountID, ".")
	for _, part := range parts {
		if !isValidAccountIDPart(part) {
			return false
		}
	}

	// Check if a 64-character account ID is a valid implicit account ID
	if len(accountID) == 64 {
		match, _ := regexp.MatchString("^[0-9a-f]+$", accountID)
		return match
	}

	return true
}

// isValidAccountIDPart checks if an account ID part follows NEAR rules
func isValidAccountIDPart(part string) bool {
	match, _ := regexp.MatchString("^[0-9a-z]+([-_][0-9a-z]+)*$", part)
	return match
}

func main() {
	// Example usage
	validAccountID := AccountID("example.near")
	invalidAccountID := AccountID("invalid.near$$$")

	// Marshaling a valid account ID
	jsonData, err := json.Marshal(validAccountID)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Valid Account ID JSON:", string(jsonData))
	}

	// Marshaling an invalid account ID
	jsonData, err = json.Marshal(invalidAccountID)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Invalid Account ID JSON:", string(jsonData))
	}

	// Unmarshaling valid JSON data into an AccountID
	jsonData = []byte(`"example.near"`)
	var accountID AccountID
	err = json.Unmarshal(jsonData, &accountID)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Unmarshaled Account ID:", accountID)
	}
}
