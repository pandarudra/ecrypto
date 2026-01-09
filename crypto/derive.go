// crypto/derive.go
package crypto

import (
    "encoding/base64"
    "errors"
    "os"
    "strings"

    "golang.org/x/crypto/argon2"
    "golang.org/x/crypto/chacha20poly1305"
)

// DeriveKeyArgon2id derives a 32-byte key from a passphrase using Argon2id.
func DeriveKeyArgon2id(pass string, salt []byte, m uint32, t uint32, p uint8) []byte {
    return argon2.IDKey(
        []byte(pass),
        salt,
        t,
        m,
        p,
        chacha20poly1305.KeySize, // 32 bytes
    )
}

// ReadKeyFromFile reads a Base64(URL)-encoded 32-byte key from a file.
func ReadKeyFromFile(path string) ([]byte, error) {
    raw, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    s := strings.TrimSpace(string(raw))

    // Try Base64URL first, then standard Base64.
    key, err := base64.RawURLEncoding.DecodeString(s)
    if err != nil {
        key, err = base64.StdEncoding.DecodeString(s)
    }
    if err != nil {
        return nil, err
    }

    if len(key) != chacha20poly1305.KeySize {
        return nil, errors.New("key must be exactly 32 bytes")
    }
    return key, nil
}

// KeySize returns the required key size (32 bytes for XChaCha20-Poly1305).
func KeySize() int {
    return chacha20poly1305.KeySize
}