package cryptostore

import "github.com/tuor4eg/ip_accounting_bot/internal/crypto"

// BaseCryptoStore provides common cryptographic key storage for any storage implementation
type BaseCryptoStore struct {
	hmacKey []byte // 32 bytes
	hmacKid int16  // e.g., 1

	aeadBox *crypto.AEADBox // AES-GCM box
	aeadKid int16           // e.g., 1
}

// SetCryptoKeys configures HMAC and AEAD encryption keys for the store
func (s *BaseCryptoStore) SetCryptoKeys(hmacKey string, hmacKid int16, aeadKey string, aeadKid int16) error {
	box, err := crypto.NewAEADBox([]byte(aeadKey))
	if err != nil {
		return err
	}

	s.hmacKey = []byte(hmacKey)
	s.hmacKid = hmacKid
	s.aeadBox = box
	s.aeadKid = aeadKid

	return nil
}

// GetHMACKey returns the HMAC key for signing
func (s *BaseCryptoStore) GetHMACKey() []byte {
	return s.hmacKey
}

// GetHMACKid returns the HMAC key ID
func (s *BaseCryptoStore) GetHMACKid() int16 {
	return s.hmacKid
}

// GetAEADBox returns the AEAD encryption box
func (s *BaseCryptoStore) GetAEADBox() *crypto.AEADBox {
	return s.aeadBox
}

// GetAEADKid returns the AEAD key ID
func (s *BaseCryptoStore) GetAEADKid() int16 {
	return s.aeadKid
}

// HasCryptoKeys checks if cryptographic keys are properly configured
func (s *BaseCryptoStore) HasCryptoKeys() bool {
	return s.hmacKey != nil && s.aeadBox != nil
}

// ExternalHash generates a hash for external ID binding to prevent cross-transport collisions
func (s *BaseCryptoStore) ExternalHash(transport, externalID string) []byte {
	// Bind the hash to transport to prevent cross-transport collisions.
	// Example: "telegram|123456789"
	payload := transport + "|" + externalID
	return crypto.HMAC256(s.hmacKey, []byte(payload))
}
