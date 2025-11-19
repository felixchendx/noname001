package sec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// credits: https://github.com/gtank/cryptopasta

func NewEncryptionKey() (string, error) {
	key := [32]byte{} // aes-256

	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key[:]), nil
}

func Encrypt(plaintext, key string) (ciphertext string, err error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("sec: malformed key: %w", err)
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", fmt.Errorf("sec: failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("sec: failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", fmt.Errorf("sec: failed to read nonce: %w", err)
	}

	cipherbytes := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(cipherbytes), nil
}

func Decrypt(ciphertext, key string) (plaintext string, err error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("sec: malformed key: %w", err)
	}
	
	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", fmt.Errorf("sec: failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("sec: failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("sec: malformed ciphertext")
	}

	cipherbytes := []byte(ciphertext)
	plainbytes, err := gcm.Open(nil, cipherbytes[:nonceSize], cipherbytes[nonceSize:], nil)
	if err != nil {
		return "", fmt.Errorf("sec: failed to decrypt ciphertext: %w", err)
	}

	return string(plainbytes), nil
}
