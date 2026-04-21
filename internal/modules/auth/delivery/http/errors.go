package http

import (
	"errors"
	"net/http"

	"github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/domain"
)

func MapDomainErrorToHTTP(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrEmptyID),
		errors.Is(err, domain.ErrEmptyName),
		errors.Is(err, domain.ErrEmptyEmail),
		errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrEmptyPassword),
		errors.Is(err, domain.ErrWeakPassword),
		errors.Is(err, domain.ErrInvalidRole):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, err.Error()

	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()

	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, err.Error()

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
