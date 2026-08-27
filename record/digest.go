package record

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"strings"
)

const digestPrefix = "sha256:"

// Digest is the normalized textual form of a SHA-256 content digest.
type Digest string

func SumBytes(data []byte) Digest {
	sum := sha256.Sum256(data)
	return Digest(digestPrefix + hex.EncodeToString(sum[:]))
}

func SumReader(r io.Reader) (Digest, int64, error) {
	h := sha256.New()
	n, err := io.Copy(h, r)
	if err != nil {
		return "", n, fmt.Errorf("hash content: %w", err)
	}
	return digestFromHash(h), n, nil
}

func DigestParts(parts ...[]byte) Digest {
	h := sha256.New()
	for _, part := range parts {
		_, _ = h.Write(part)
	}
	return digestFromHash(h)
}

func ParseDigest(raw string) (Digest, error) {
	d := Digest(strings.ToLower(strings.TrimSpace(raw)))
	if err := d.Validate(); err != nil {
		return "", err
	}
	return d, nil
}

func (d Digest) Validate() error {
	raw := string(d)
	if !strings.HasPrefix(raw, digestPrefix) {
		return fmt.Errorf("digest %q must start with %q", raw, digestPrefix)
	}
	hexPart := strings.TrimPrefix(raw, digestPrefix)
	if len(hexPart) != sha256.Size*2 {
		return fmt.Errorf("digest %q has %d hex characters; want %d", raw, len(hexPart), sha256.Size*2)
	}
	decoded, err := hex.DecodeString(hexPart)
	if err != nil {
		return fmt.Errorf("digest %q contains invalid hex: %w", raw, err)
	}
	if len(decoded) != sha256.Size {
		return errors.New("decoded digest has unexpected length")
	}
	return nil
}

func (d Digest) Hex() (string, error) {
	if err := d.Validate(); err != nil {
		return "", err
	}
	return strings.TrimPrefix(string(d), digestPrefix), nil
}

func (d Digest) String() string { return string(d) }

func digestFromHash(h hash.Hash) Digest {
	return Digest(digestPrefix + hex.EncodeToString(h.Sum(nil)))
}
