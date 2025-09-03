package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"

	"github.com/tuor4eg/ip_accounting_bot/internal/validate"
)

// NewAEADBox creates an AES-GCM instance using the provided 32-byte key.
// AES-GCM = authenticated encryption: ensures both confidentiality and integrity.
func NewAEADBox(key []byte) (*AEADBox, error) {
	const op = "crypto.NewAEADBox"

	if len(key) != 32 {
		return nil, validate.Wrap(op, ErrInvalidKey)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, validate.Wrap(op, err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, validate.Wrap(op, err)
	}

	return &AEADBox{aead: aead}, nil
}

// Seal encrypts the plaintext using AES-GCM.
// - `aad` = Additional Authenticated Data: binds the ciphertext to a context (e.g., "telegram").
// - The output format: [1-byte version][nonce][ciphertext+tag].
func (b *AEADBox) Seal(plaintext, aad []byte) ([]byte, error) {
	const op = "crypto.AEADBox.Seal"

	nonce := make([]byte, b.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, validate.Wrap(op, err)
	}

	// Prefix with version = 1 for future compatibility (in case we change algorithms).
	out := make([]byte, 1+len(nonce))
	out[0] = 1
	copy(out[1:], nonce)

	ciphertext := b.aead.Seal(nil, nonce, plaintext, aad)
	return append(out, ciphertext...), nil
}

// Open decrypts the ciphertext produced by Seal.
// If the ciphertext, nonce, AAD, or tag is invalid, this returns an error.
func (b *AEADBox) Open(box, aad []byte) ([]byte, error) {
	const op = "crypto.AEADBox.Open"

	ns := b.aead.NonceSize()
	if len(box) < 1+ns {
		return nil, validate.Wrap(op, ErrCipherTooShort)
	}

	nonce := box[1 : 1+ns]
	ct := box[1+ns:]

	return b.aead.Open(nil, nonce, ct, aad)
}

// HMAC256 calculates HMAC-SHA256(data) using the provided key.
// Used for securely deriving external hashes without storing raw IDs.
func HMAC256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// PackInt64 converts an int64 into 8 bytes (BigEndian).
func PackInt64(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// UnpackInt64 converts 8 bytes back into int64.
func UnpackInt64(b []byte) (int64, error) {
	const op = "crypto.UnpackInt64"

	if len(b) != 8 {
		return 0, validate.Wrap(op, ErrInvalidInt64Length)
	}
	return int64(binary.BigEndian.Uint64(b)), nil
}

// EncryptInt64 encrypts an int64 value directly with AES-GCM.
func EncryptInt64(box *AEADBox, v int64, aad []byte) ([]byte, error) {
	return box.Seal(PackInt64(v), aad)
}

// DecryptInt64 decrypts an AES-GCM ciphertext back into an int64.
func DecryptInt64(box *AEADBox, enc []byte, aad []byte) (int64, error) {
	const op = "crypto.DecryptInt64"

	pt, err := box.Open(enc, aad)
	if err != nil {
		return 0, validate.Wrap(op, err)
	}
	return UnpackInt64(pt)
}
