package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// encryption ref
// https://github.com/phachon/fasthttpsession/blob/master/encrypt.go
// https://github.com/gofiber/fiber/blob/main/middleware/encryptcookie/encryptcookie.go

const (
	// TODO: move to web config
	// TODO: add support for dynamic key rotation without breaking session
	veryRandomKey = "UBBtmvF30483auDiCDatE4DgOYUrBvDFhxB1JGhN3xU="
)

func (store *CookieStore) GenerateKey() ([]byte, error) {
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return []byte(""), err
	}

	genKey := base64.StdEncoding.EncodeToString(key)

	return []byte(genKey), nil
}

// TODO: move to package sec as general use enc dec
func (store *CookieStore) EncryptCookie(val []byte) ([]byte, error) {
	emptyByte := []byte("")

	keyDecoded, err := base64.StdEncoding.DecodeString(veryRandomKey)
	if err != nil {
		return emptyByte, fmt.Errorf("failed to base64-decode key: %w", err)
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		return emptyByte, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return emptyByte, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return emptyByte, fmt.Errorf("failed to read nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, val, nil)

	encodedVal := base64.StdEncoding.EncodeToString(ciphertext)

	return []byte(encodedVal), nil
}

func (store *CookieStore) DecryptCookie(val []byte) ([]byte, error) {
	emptyByte := []byte("")

	keyDecoded, err := base64.StdEncoding.DecodeString(veryRandomKey)
	if err != nil {
		return emptyByte, fmt.Errorf("failed to base64-decode key: %w", err)
	}

	enc, err := base64.StdEncoding.DecodeString(string(val))
	if err != nil {
		return emptyByte, fmt.Errorf("failed to base64-decode value: %w", err)
	}

	block, err := aes.NewCipher(keyDecoded)
	if err != nil {
		return emptyByte, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return emptyByte, fmt.Errorf("failed to create GCM mode: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(enc) < nonceSize {
		return emptyByte, fmt.Errorf("invalid encrypted value")
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return emptyByte, fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	return plaintext, nil
}

// type CookieConfig struct {
// 	CookieName string
// 	CookieDomain string
// 	CookiePath string
// 	// CookieSameSite string
// 	Expiration int
// 	CookieSecure bool
// 	CookieHTTPOnly bool
// 	// CookieSessionOnly bool
// }
