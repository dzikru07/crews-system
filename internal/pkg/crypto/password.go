package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    uint32 = 3
	argon2Memory  uint32 = 64 * 1024
	argon2Threads uint8  = 2
	argon2KeyLen  uint32 = 32
	argon2SaltLen        = 16
)

// HashSecret creates an Argon2id hash from a secret.
//
// The returned value contains all parameters and the random salt
// required to verify the secret later.
//
// Format:
//
// $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
func HashSecret(secret string) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("secret cannot be empty")
	}

	salt := make([]byte, argon2SaltLen)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(secret),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argon2Memory,
		argon2Time,
		argon2Threads,
		saltEncoded,
		hashEncoded,
	), nil
}

// VerifySecret verifies a plaintext secret against an Argon2id hash.
func VerifySecret(secret, encodedHash string) bool {
	if secret == "" || encodedHash == "" {
		return false
	}

	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return false
	}

	if parts[1] != "argon2id" {
		return false
	}

	if parts[2] != "v=19" {
		return false
	}

	var memory uint32
	var time uint32
	var threads uint8

	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&time,
		&threads,
	); err != nil {
		return false
	}

	if memory == 0 || time == 0 || threads == 0 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	if len(salt) == 0 || len(expectedHash) == 0 {
		return false
	}

	actualHash := argon2.IDKey(
		[]byte(secret),
		salt,
		time,
		memory,
		threads,
		uint32(len(expectedHash)),
	)

	if len(actualHash) != len(expectedHash) {
		return false
	}

	return subtle.ConstantTimeCompare(
		actualHash,
		expectedHash,
	) == 1
}
