package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/estenssoros/relay/config"
	"github.com/pkg/errors"
)

func newGCM() (cipher.AEAD, error) {
	key, err := config.CipherKeyBytes()
	if err != nil {
		return nil, errors.Wrap(err, "cypher key bytes")
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	return gcm, nil
}

// Decrypt translate an encoded string
func Decrypt(in string) (string, error) {
	gcm, err := newGCM()
	if err != nil {
		return "", errors.Wrap(err, "new gcm")
	}
	cipherBytes, err := hex.DecodeString(in)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherBytes) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	gcmBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(gcmBytes), nil

}

// Encrypt encrypt  a string
func Encrypt(in string) (string, error) {
	gcm, err := newGCM()
	if err != nil {
		return "", errors.Wrap(err, "new gcm")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", gcm.Seal(nonce, nonce, []byte(in), nil)), nil
}
