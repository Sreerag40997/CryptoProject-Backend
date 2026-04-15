package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt(text string, key []byte) (string, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(text), nil)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(enc string, key []byte) (string, error) {

	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()

	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}