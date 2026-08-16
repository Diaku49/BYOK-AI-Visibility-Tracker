package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type KeyCipher struct {
	aead cipher.AEAD
}

func NewKeyCipher(masterKey []byte) (*KeyCipher, error) {
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, fmt.Errorf("create aes cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	return &KeyCipher{
		aead: aead,
	}, nil
}

func NewKeyCipherFromBase64(encodedMasterKey string) (*KeyCipher, error) {
	masterKey, err := base64.StdEncoding.DecodeString(encodedMasterKey)
	if err != nil {
		return nil, fmt.Errorf("decode master key: %w", err)
	}

	return NewKeyCipher(masterKey)
}

func (c *KeyCipher) Encrypt(plaintext []byte) (ciphertext []byte, nonce []byte, err error) {
	nonce = make([]byte, c.aead.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext = c.aead.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

func (c *KeyCipher) Decrypt(ciphertext []byte, nonce []byte) ([]byte, error) {
	plaintext, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt provider key: %w", err)
	}

	return plaintext, nil
}
