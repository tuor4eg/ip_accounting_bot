package cryptostore

import "github.com/tuor4eg/ip_accounting_bot/internal/crypto"

// CryptoStore defines cryptographic key storage capabilities that any storage can implement
type CryptoStore interface {
	SetCryptoKeys(hmacKey string, hmacKid int16, aeadKey string, aeadKid int16) error
	GetHMACKey() []byte
	GetHMACKid() int16
	GetAEADBox() *crypto.AEADBox
	GetAEADKid() int16
	HasCryptoKeys() bool
	ExternalHash(transport, externalID string) []byte
}
