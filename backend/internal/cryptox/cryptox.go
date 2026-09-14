// Package cryptox handles encryption of secrets (Agnes API key) at rest
// using AES-256-GCM with a key derived from a master key file.
package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

// Cipher is a small wrapper around AES-256-GCM with a 32-byte key.
type Cipher struct {
	gcm cipher.AEAD
}

func New(key []byte) (*Cipher, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("cryptox: key must be exactly 32 bytes, got %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Cipher{gcm: g}, nil
}

// NewFromBase64 loads a key from a base64 string.
func NewFromBase64(s string) (*Cipher, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return New(b)
}

// NewFromFile loads (or creates) a key file at path. If the file does not
// exist a fresh 32-byte key is generated and persisted (0600).
func NewFromFile(path string) (*Cipher, error) {
	if data, err := os.ReadFile(path); err == nil {
		b, err := base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			return nil, fmt.Errorf("decode master key file: %w", err)
		}
		return New(b)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	// generate
	k := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil, err
	}
	encoded := base64.StdEncoding.EncodeToString(k)
	if err := os.WriteFile(path, []byte(encoded), 0o600); err != nil {
		return nil, err
	}
	return New(k)
}

// Derive derives a 32-byte key from a passphrase using SHA-256. This is
// intentionally simple — meant only for the AGNES_STUDIO_MASTER_KEY env
// var path, not the file path.
func Derive(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

func (c *Cipher) Encrypt(plaintext []byte) (string, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := c.gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ct), nil
}

func (c *Cipher) Decrypt(ciphertext string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	if len(raw) < c.gcm.NonceSize() {
		return nil, errors.New("cryptox: ciphertext too short")
	}
	nonce := raw[:c.gcm.NonceSize()]
	ct := raw[c.gcm.NonceSize():]
	return c.gcm.Open(nil, nonce, ct, nil)
}
