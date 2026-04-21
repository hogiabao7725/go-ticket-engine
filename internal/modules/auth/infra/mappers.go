package infra

import (
	"github.com/hogiabao7725/gin-auth-playground/internal/infra/sqlc"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain/uservo"
)

// toDomainUser: sqlc.User -> domain.User
func toDomainUser(dbUser *sqlc.User) *domain.User {
	name := uservo.ReconstituteName(dbUser.Name)
	email := uservo.ReconstituteEmail(dbUser.Email)
	hashedPwd := uservo.ReconstituteHashedPassword(dbUser.Password)

	role := uservo.ReconstituteRole(dbUser.Role)

	return domain.ReconstructUser(
		dbUser.ID,
		name,
		email,
		hashedPwd,
		role,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	)
}
