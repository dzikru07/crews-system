package credential

import (
	"testing"

	"github.com/dzikru07/crews-system/internal/pkg/crypto"
)

func TestCredentialSecretVerification(t *testing.T) {
	secret := "test-secret-123"

	hash, err := crypto.HashSecret(secret)
	if err != nil {
		t.Fatalf("failed to hash secret: %v", err)
	}

	if !crypto.VerifySecret(secret, hash) {
		t.Fatal("expected secret to be valid")
	}

	if crypto.VerifySecret("wrong-secret", hash) {
		t.Fatal("expected wrong secret to be invalid")
	}
}
