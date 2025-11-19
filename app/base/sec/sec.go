package sec

import (
	"fmt"
)

type DumbSecurityBundle struct {
	keyring map[string]string
}
func NewDumbSecurityBundle() (*DumbSecurityBundle) {
	return &DumbSecurityBundle{
		keyring: make(map[string]string),
	}
}

func (secBundle *DumbSecurityBundle) AddKey(label, key string) {
	secBundle.keyring[label] = key
}

func (secBundle *DumbSecurityBundle) Encrypt(plaintext, keyLabel string) (ciphertext string, err error) {
	key, found := secBundle.keyring[keyLabel]
	if !found {
		return "", fmt.Errorf("sec: key[%s] not found.", keyLabel)
	}

	return Encrypt(plaintext, key)
}

func (secBundle *DumbSecurityBundle) Decrypt(ciphertext, keyLabel string) (plaintext string, err error) {
	key, found := secBundle.keyring[keyLabel]
	if !found {
		return "", fmt.Errorf("sec: key[%s] not found.", keyLabel)
	}

	return Decrypt(ciphertext, key)
}
