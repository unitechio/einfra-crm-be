package algorithm

import (
	"os"
	"testing"
)

func TestEncryptDecryptFile(t *testing.T) {
	plainFile := "test_plain.txt"
	encFile := "test_cipher.txt"
	decFile := "test_decrypted.txt"

	// Tạo file plaintext mẫu
	plainText := []byte("Đây là nội dung test file AES-192 CBC PKCS7 với IV ngẫu nhiên.")
	if err := os.WriteFile(plainFile, plainText, 0644); err != nil {
		t.Fatalf("Tạo file plaintext thất bại: %v", err)
	}
	defer os.Remove(plainFile)
	defer os.Remove(encFile)
	defer os.Remove(decFile)

	// Mã hóa file
	if err := EncryptFile(plainFile, encFile); err != nil {
		t.Fatalf("EncryptFile thất bại: %v", err)
	}

	// Giải mã file
	if err := DecryptFile(encFile, decFile); err != nil {
		t.Fatalf("DecryptFile thất bại: %v", err)
	}

	// So sánh nội dung
	decData, err := os.ReadFile(decFile)
	if err != nil {
		t.Fatalf("Đọc file decrypted thất bại: %v", err)
	}

	if string(decData) != string(plainText) {
		t.Fatalf("Round-trip file thất bại. Expected: %s, Got: %s", plainText, decData)
	}
}
