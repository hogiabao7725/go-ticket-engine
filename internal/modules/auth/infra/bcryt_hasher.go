package infra

import (
	"github.com/hogiabao7725/go-ticket-engine/internal/infra/crypto"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/domain"
)

type BcryptHasher struct {
	bcrypt crypto.Bcrypt
}

var _ domain.PasswordHasher = (*BcryptHasher)(nil)

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(plain string) (string, error) {
	return h.bcrypt.Hash(plain)
}

func (h *BcryptHasher) Compare(hash, plain string) error {
	return h.bcrypt.Compare(hash, plain)
}
