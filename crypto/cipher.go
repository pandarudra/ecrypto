// crypto/cipher.go
package crypto

import (
    "errors"

    "golang.org/x/crypto/chacha20poly1305"
)

// EncryptAEAD encrypts plaintext with XChaCha20-Poly1305.
// aad: additional authenticated data (header).
// nonce: 24-byte nonce (must be random).
// Returns ciphertext with authentication tag appended.
func EncryptAEAD(key, plaintext, aad, nonce []byte) ([]byte, error) {
    if len(key) != chacha20poly1305.KeySize {
        return nil, errors.New("key must be 32 bytes")
    }
    if len(nonce) != 24 {
        return nil, errors.New("nonce must be 24 bytes")
    }

    aead, err := chacha20poly1305.NewX(key)
    if err != nil {
        return nil, err
    }

    ciphertext := aead.Seal(nil, nonce, plaintext, aad)
    return ciphertext, nil
}

// DecryptAEAD decrypts ciphertext with XChaCha20-Poly1305.
// aad: additional authenticated data (header).
// nonce: 24-byte nonce (must match encryption nonce).
// Returns plaintext or error if authentication fails.
func DecryptAEAD(key, ciphertext, aad, nonce []byte) ([]byte, error) {
    if len(key) != chacha20poly1305.KeySize {
        return nil, errors.New("key must be 32 bytes")
    }
    if len(nonce) != 24 {
        return nil, errors.New("nonce must be 24 bytes")
    }

    aead, err := chacha20poly1305.NewX(key)
    if err != nil {
        return nil, err
    }

    plaintext, err := aead.Open(nil, nonce, ciphertext, aad)
    if err != nil {
        return nil, errors.New("decryption failed: authentication tag mismatch or wrong key")
    }
    return plaintext, nil
}