package infra

import (
	"fmt"

	infraToken "github.com/hogiabao7725/gin-auth-playground/internal/infra/token"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
)

type JWTTokenGenerator struct {
	jwt *infraToken.JWT
}

var _ domain.TokenGenerator = (*JWTTokenGenerator)(nil)

func NewJWTTokenGenerator(jwt *infraToken.JWT) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		jwt: jwt,
	}
}

func (g *JWTTokenGenerator) GenerateAccessToken(userID, role string) (domain.TokenResult, error) {
	token, err := g.jwt.GenerateAccessToken(userID, role)
	if err != nil {
		return domain.TokenResult{}, fmt.Errorf("auth.infra.token_generator.GenerateAccessToken: %w", err)
	}
	return domain.TokenResult{
		Token:     token,
		ExpiresIn: g.jwt.AccessTTL(),
	}, nil
}

func (g *JWTTokenGenerator) GenerateRefreshToken(userID string) (domain.TokenResult, error) {
	token, err := g.jwt.GenerateRefreshToken(userID)
	if err != nil {
		return domain.TokenResult{}, fmt.Errorf("auth.infra.token_generator.GenerateRefreshToken: %w", err)
	}
	return domain.TokenResult{
		Token:     token,
		ExpiresIn: g.jwt.RefreshTTL(),
	}, nil
}
