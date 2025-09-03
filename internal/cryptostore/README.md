# CryptoStore Package

This package provides cryptographic key storage capabilities that can be embedded into storage implementations.

## Purpose

The `BaseCryptoStore` is designed to store cryptographic keys only. The actual encryption/decryption operations are handled by the `crypto` package.

## Components

### BaseCryptoStore

Base structure providing cryptographic key storage:
- HMAC key storage for data signing
- AEAD key storage for encryption operations
- Key versioning support (key IDs)

## Usage

### Embedded Usage

```go
import "github.com/tuor4eg/ip_accounting_bot/internal/cryptostore"

type MyStorage struct {
    cryptostore.BaseCryptoStore // Embed crypto key storage
    // ... other storage fields
}

// MyStorage automatically gets all crypto key methods
```

### Key Management

```go
// Configure keys
err := store.SetCryptoKeys("hmac-key", 1, "aead-key", 1)
if err != nil {
    // Handle error
}

// Access keys when needed
if store.HasCryptoKeys() {
    hmacKey := store.GetHMACKey()
    aeadKey := store.GetAEADKey()
    
    // Use keys with crypto package
    // crypto.Sign(hmacKey, data)
    // crypto.Encrypt(aeadKey, data)
}
```

## Interface

All implementations provide the `CryptoStore` interface:

```go
type CryptoStore interface {
    SetCryptoKeys(hmacKey string, hmacKid int16, aeadKey string, aeadKid int16) error
    GetHMACKey() []byte
    GetHMACKid() int16
    GetAEADKey() []byte
    GetAEADKid() int16
    HasCryptoKeys() bool
}
```

## Architecture

- **Key Storage**: `BaseCryptoStore` stores the keys
- **Encryption Logic**: `crypto` package handles actual encryption/decryption
- **Separation of Concerns**: Storage stores keys, crypto package uses them

This design allows storage implementations to focus on data persistence while providing access to cryptographic keys when needed by the crypto package.
