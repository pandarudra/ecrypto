// crypto/types.go
package crypto

import (
    "bytes"
    "encoding/binary"
    "errors"
    "io"
)

// HeaderV1 is the .ecrypt container header (v1 format).
type HeaderV1 struct {
    Magic   [8]byte // "ECRYPT01"
    Version uint8   // 1
    KDF     uint8   // 0=raw key, 1=Argon2id
    ArgonM  uint32  // Argon2 memory (KiB)
    ArgonT  uint32  // Argon2 time cost
    ArgonP  uint8   // Argon2 parallelism
    Salt    [16]byte
    Nonce   [24]byte
}

// Encode serializes HeaderV1 to bytes.
func (h *HeaderV1) Encode() []byte {
    var buf bytes.Buffer
    _ = binary.Write(&buf, binary.LittleEndian, h.Magic)
    _ = binary.Write(&buf, binary.LittleEndian, h.Version)
    _ = binary.Write(&buf, binary.LittleEndian, h.KDF)
    _ = binary.Write(&buf, binary.LittleEndian, h.ArgonM)
    _ = binary.Write(&buf, binary.LittleEndian, h.ArgonT)
    _ = binary.Write(&buf, binary.LittleEndian, h.ArgonP)
    _ = binary.Write(&buf, binary.LittleEndian, h.Salt)
    _ = binary.Write(&buf, binary.LittleEndian, h.Nonce)
    return buf.Bytes()
}

// DecodeHeaderV1 deserializes HeaderV1 from a reader.
func DecodeHeaderV1(r io.Reader) (*HeaderV1, error) {
    h := &HeaderV1{}
    if err := binary.Read(r, binary.LittleEndian, &h.Magic); err != nil {
        return nil, err
    }
    if string(h.Magic[:]) != "ECRYPT01" {
        return nil, errors.New("invalid magic: not an ecrypto container")
    }
    if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
        return nil, err
    }
    if h.Version != 1 {
        return nil, errors.New("unsupported container version")
    }
    if err := binary.Read(r, binary.LittleEndian, &h.KDF); err != nil {
        return nil, err
    }
    if err := binary.Read(r, binary.LittleEndian, &h.ArgonM); err != nil {
        return nil, err
    }
    if err := binary.Read(r, binary.LittleEndian, &h.ArgonT); err != nil {
        return nil, err
    }
    if err := binary.Read(r, binary.LittleEndian, &h.ArgonP); err != nil {
        return nil, err
    }
    if err := binary.Read(r, binary.LittleEndian, &h.Salt); err != nil {
        return nil, err
    }
    if err := binary.Read(r, binary.LittleEndian, &h.Nonce); err != nil {
        return nil, err
    }
    return h, nil
}

// HeaderSize returns the byte size of an encoded HeaderV1.
func HeaderSize() int {
    return 8 + 1 + 1 + 4 + 4 + 1 + 16 + 24 // 59 bytes
}