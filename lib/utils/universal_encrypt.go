package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	//"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	//"encoding/json"
	"errors"
	//"fmt"
	"io"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// MasterSecret should be loaded from an environment variable for production
var MasterSecret = []byte("a-very-secret-32-byte-long-key!!")

type Crypter struct {
	Key []byte
}

// NewCrypter initializes the struct with the default master key
func NewCrypter(key []byte) *Crypter {
	return &Crypter{Key: key}
}

// ============================================================================
// MODERN AUTHENTICATED ENCRYPTION (AES-GCM) - Used for API Keys & Secrets
// ============================================================================

// EncryptSecret encrypts a string into a Base64 string using AES-GCM
func (c *Crypter) EncryptSecret(plaintext string) (string, error) {
	block, err := aes.NewCipher(c.Key)
	if err != nil { return "", err }

	gcm, err := cipher.NewGCM(block)
	if err != nil { return "", err }

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil { return "", err }

	// Seal appends the ciphertext to the nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSecret decodes a Base64 string and decrypts it using AES-GCM
func (c *Crypter) DecryptSecret(cryptoText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil { return "", err }

	block, err := aes.NewCipher(c.Key)
	if err != nil { return "", err }

	gcm, err := cipher.NewGCM(block)
	if err != nil { return "", err }

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize { return "", errors.New("ciphertext too short") }

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil { return "", err }

	return string(plaintext), nil
}

// ============================================================================
// PASSWORD & UTILITY HASHING
// ============================================================================

func (c *Crypter) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func (c *Crypter) CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (c *Crypter) Md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ============================================================================
// BASE64 HELPERS
// ============================================================================

func (c *Crypter) Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func (c *Crypter) Base64Decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(data)
}