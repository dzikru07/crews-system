package crypto

import "testing"

func TestHashSecret(t *testing.T) {
	secret := "crs_live_super_secret"

	hash, err := HashSecret(secret)
	if err != nil {
		t.Fatalf("HashSecret() error = %v", err)
	}

	if hash == "" {
		t.Fatal("expected hash, got empty string")
	}

	if hash == secret {
		t.Fatal("secret must not be stored as plaintext")
	}

	if len(hash) < 50 {
		t.Fatalf("hash looks too short: %d", len(hash))
	}

	t.Logf("generated hash: %s", hash)
}

func TestVerifySecret(t *testing.T) {
	secret := "crs_live_super_secret"

	hash, err := HashSecret(secret)
	if err != nil {
		t.Fatalf("HashSecret() error = %v", err)
	}

	if !VerifySecret(secret, hash) {
		t.Fatal("expected valid secret to verify")
	}

	if VerifySecret("wrong-secret", hash) {
		t.Fatal("expected wrong secret to fail verification")
	}
}

func TestHashSecretGeneratesDifferentHashes(t *testing.T) {
	secret := "same-secret"

	hash1, err := HashSecret(secret)
	if err != nil {
		t.Fatalf("HashSecret() error = %v", err)
	}

	hash2, err := HashSecret(secret)
	if err != nil {
		t.Fatalf("HashSecret() error = %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("expected different hashes because salts should be random")
	}

	if !VerifySecret(secret, hash1) {
		t.Fatal("hash1 should verify")
	}

	if !VerifySecret(secret, hash2) {
		t.Fatal("hash2 should verify")
	}
}

func TestVerifySecretInvalidHash(t *testing.T) {
	testCases := []string{
		"",
		"invalid",
		"$argon2id$invalid",
		"$argon2id$v=19$m=65536,t=3,p=2$invalid$invalid",
	}

	for _, hash := range testCases {
		if VerifySecret("secret", hash) {
			t.Errorf("expected invalid hash to fail: %q", hash)
		}
	}
}

func TestHashSecretEmptySecret(t *testing.T) {
	_, err := HashSecret("")

	if err == nil {
		t.Fatal("expected error for empty secret")
	}
}
