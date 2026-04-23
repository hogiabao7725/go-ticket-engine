package infra

import (
	"github.com/hogiabao7725/gin-auth-playground/internal/infra/crypto"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
)

type SHA256Hasher struct {
	sha256 crypto.SHA256
}

var _ domain.TokenHasher = (*SHA256Hasher)(nil)

func NewSHA256Hasher() *SHA256Hasher {
	return &SHA256Hasher{}
}

func (h *SHA256Hasher) Hash(plain string) (string, error) {
	return h.sha256.Hash(plain)
}

func (h *SHA256Hasher) Compare(hash, plain string) error {
	return h.sha256.Compare(hash, plain)
}
