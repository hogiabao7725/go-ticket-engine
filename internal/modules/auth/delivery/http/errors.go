package http

import (
	"errors"
	"net/http"

	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/domain"
)

func mapDomainErrorToHTTP(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, domain.ErrWeakPassword),
		errors.Is(err, domain.ErrInvalidName),
		errors.Is(err, domain.ErrInvalidEmail):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
