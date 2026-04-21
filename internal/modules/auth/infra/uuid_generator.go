package infra

import (
	"github.com/hogiabao7725/gin-auth-playground/internal/infra/identifier"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
)

type UUIDGenerator struct {
	gen identifier.UUID
}

var _ domain.IdentifierGenerator = (*UUIDGenerator)(nil)

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) Generate() string {
	return g.gen.Generate()
}
