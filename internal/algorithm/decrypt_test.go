package algorithm

import (
	"testing"
)

func TestRoundTrip(t *testing.T) {
	original := "Some secret text to test AES-192"
	cipherText, err := Encrypt([]byte(original))
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(cipherText)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != original {
		t.Fatalf("Round-trip failed. Expected: %s, Got: %s", original, decrypted)
	}
}
