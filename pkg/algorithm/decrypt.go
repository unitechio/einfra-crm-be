package algorithm

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

// PKCS7 unpad
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}
	padding := int(data[length-1])
	if padding > length || padding == 0 {
		return nil, errors.New("invalid padding")
	}
	return data[:length-padding], nil
}

// Decrypt chuỗi base64
func Decrypt(cipherTextB64 string) ([]byte, error) {
	cipherWithIV, err := base64.StdEncoding.DecodeString(cipherTextB64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherWithIV) < block.BlockSize() {
		return nil, errors.New("ciphertext too short")
	}

	iv := cipherWithIV[:block.BlockSize()]
	cipherText := cipherWithIV[block.BlockSize():]

	if len(cipherText)%block.BlockSize() != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	return pkcs7Unpad(plainText)
}

// DecryptFile đọc file input (base64), giải mã, lưu ra output
func DecryptFile(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read input file failed: %w", err)
	}

	dec, err := Decrypt(string(data))
	if err != nil {
		return fmt.Errorf("decrypt failed: %w", err)
	}

	return os.WriteFile(outputPath, dec, 0644)
}
