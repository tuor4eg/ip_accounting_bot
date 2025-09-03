package cryptostore

import "github.com/tuor4eg/ip_accounting_bot/internal/crypto"

// BaseCryptoStore provides common cryptographic key storage for any storage implementation
type BaseCryptoStore struct {
	hmacKey []byte // 32 bytes
	hmacKid int16  // e.g., 1

	aeadBox *crypto.AEADBox // AES-GCM box
	aeadKid int16           // e.g., 1
}
