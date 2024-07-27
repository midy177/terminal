package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

const key = "hkqspsimisnj57ps"

func AesEncryptByGCM(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceStr := key[:gcm.NonceSize()]
	nonce := []byte(nonceStr)
	seal := gcm.Seal(nonce, nonce, data, nil)
	return seal, nil
}

func AesDecryptByGCM(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	open, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return open, nil
}
