package near

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/OCD-Labs/store-hub/util"
)

type credentials struct {
	AcctID  string `json:"account_id"`
	PubKey  string `json:"public_key"`
	PrivKey string `json:"private_key"`
}

func InstallNearCLI() (err error) {
	if !util.CommandExists("near") {
		cmd := exec.Command("npm", "install", "-g", "near-cli")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			return fmt.Errorf("failed to install near-cli: %v", err)
		}
	}
	return nil
}

func RunNearCLICommand(args ...string) error {
	cmd := exec.Command("near", args...)
	cmd.Env = append(os.Environ(), "NEAR_ENV=testnet")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return fmt.Errorf("failed to run near-cli command: %v", err)
	}
	fmt.Println(string(output))
	return nil
}

func SetupNearMasterAccount(network, accountID, pubKey, privKey string) (err error) {
	if !strings.HasSuffix(accountID, ".near") && !strings.HasSuffix(accountID, ".testnet") {
		return fmt.Errorf("accout_id must contain '.near' or '.testnet'")
	}

	if network != "testnet" && network != "mainnet" {
		return fmt.Errorf("network must be mainnet or testnet")
	}

	if !strings.HasPrefix(pubKey, "ed25519:") || !strings.HasPrefix(privKey, "ed25519:") {
		return fmt.Errorf("public and private keys must start with 'ed25519:'")
	}

	homePath := os.Getenv("HOME")
	credentialsFolder := filepath.Join(homePath, ".near-credentials", network)
	if !util.FolderExists(credentialsFolder) {
		if err = os.MkdirAll(credentialsFolder, os.ModePerm); err != nil {
			return err
		}
	}

	credentialsPath := filepath.Join(credentialsFolder, fmt.Sprintf("%s.json", accountID))
	if !util.FolderExists(credentialsPath) {
		acctCred := credentials{
			AcctID:  accountID,
			PubKey:  pubKey,
			PrivKey: privKey,
		}
		err = writeJSONToFile(acctCred, credentialsPath)
		if err != nil {
			return fmt.Errorf("failed to write JSON to file: %w", err)
		}
	}

	return nil
}

func writeJSONToFile(data interface{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}
