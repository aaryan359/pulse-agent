package agent

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type ServerIdentity struct {
	ServerID   string `json:"server_id"`
	APIKeyHash string `json:"api_key_hash"`
}

func HashAPIKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

func identityPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pulse", "identity.json")
}

func SaveServerIdentity(serverID, apiKey string) error {
	identity := ServerIdentity{
		ServerID:   serverID,
		APIKeyHash: HashAPIKey(apiKey),
	}

	data, _ := json.Marshal(identity)
	_ = os.MkdirAll(filepath.Dir(identityPath()), 0700)
	return os.WriteFile(identityPath(), data, 0600)
}

func LoadServerIdentity() (*ServerIdentity, error) {
	data, err := os.ReadFile(identityPath())
	if err != nil {
		return nil, err
	}

	var identity ServerIdentity
	if err := json.Unmarshal(data, &identity); err != nil {
		return nil, err
	}

	return &identity, nil
}

func ClearServerIdentity() {
	_ = os.Remove(identityPath())
}
