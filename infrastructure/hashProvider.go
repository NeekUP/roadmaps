package infrastructure

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"github.com/NeekUP/roadmaps/core"
)

type Sha256HashProvider struct {
}

func NewSha256HashProvider() core.HashProvider {
	return &Sha256HashProvider{}
}

func (provider *Sha256HashProvider) HashPassword(pass string) (hash []byte, salt []byte) {
	s := make([]byte, 32)
	rand.Read(s)

	h := sha256.Sum256(append([]byte(pass), s...))
	return h[:], s
}

func (provider *Sha256HashProvider) CheckPassword(pass string, hash []byte, salt []byte) bool {
	h := sha256.Sum256(append([]byte(pass), salt...))
	return bytes.Compare(hash, h[:]) == 0
}
