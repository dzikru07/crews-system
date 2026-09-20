package credential

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateClientID(prefix string) (string, error) {
	token, err := randomToken(24)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s_%s", prefix, token), nil
}

func GenerateClientSecret(prefix string) (string, error) {
	token, err := randomToken(48)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s_%s", prefix, token), nil
}

func randomToken(size int) (string, error) {
	data := make([]byte, size)

	if _, err := rand.Read(data); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}
