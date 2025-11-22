package algorithm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// AES-192 key (24 bytes) hardcoded mẫu
var key = []byte{161, 227, 194, 25, 25, 173, 238, 171, 12, 34, 56, 78, 90, 123, 45, 67, 89, 101, 112, 131, 145, 167, 189, 210}

// PKCS7 padding
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// Encrypt chuỗi và trả về base64
func Encrypt(plainText []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainText = pkcs7Pad(plainText, block.BlockSize())

	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	// Gắn IV tiền tố vào ciphertext
	cipherWithIV := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(cipherWithIV), nil
}

// EncryptFile đọc file input, mã hóa, lưu ra output
func EncryptFile(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file failed: %w", err)
	}

	enc, err := Encrypt(data)
	if err != nil {
		return fmt.Errorf("encrypt failed: %w", err)
	}

	return os.WriteFile(outputPath, []byte(enc), 0644)
}
