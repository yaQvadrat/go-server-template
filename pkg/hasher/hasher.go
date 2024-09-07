package hasher

import (
	"crypto/sha1"
	"fmt"
)

type Hasher interface {
	Hash(input string) string
}

type SHA1Hasher struct {
	HashSalt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt}
}

func (s *SHA1Hasher) Hash(input string) string {
	hash := sha1.New()
	hash.Write([]byte(input))

	return fmt.Sprintf("%x", hash.Sum([]byte(s.HashSalt)))
}
