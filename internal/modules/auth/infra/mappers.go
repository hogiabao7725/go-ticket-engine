package infra

import (
	"github.com/hogiabao7725/go-ticket-engine/internal/infra/sqlc"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/domain"
)

// toDomainUser: sqlc.User -> domain.User
func toDomainUser(dbUser *sqlc.User) *domain.User {
	name := domain.ReconstituteName(dbUser.Name)
	email := domain.ReconstituteEmail(dbUser.Email)
	hashedPwd := domain.ReconstituteHashedPassword(dbUser.Password)

	role := domain.ReconstituteRole(dbUser.Role)

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
