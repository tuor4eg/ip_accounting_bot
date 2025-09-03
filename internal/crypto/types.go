package crypto

import "crypto/cipher"

// AEADBox provides authenticated encryption using AES-GCM
type AEADBox struct {
	aead cipher.AEAD
}
