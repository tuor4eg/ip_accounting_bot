package crypto

import "errors"

var (
	ErrInvalidKey         = errors.New("key must be 32 bytes")
	ErrCipherTooShort     = errors.New("ciphertext too short")
	ErrInvalidInt64Length = errors.New("invalid int64 length")
)
