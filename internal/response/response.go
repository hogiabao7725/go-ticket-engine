package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hogiabao7725/go-ticket-engine/internal/apperror"
)

type SuccessResponse struct {
	Data any `json:"data"`
}

type ValidationDetail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(
		http.StatusOK,
		&SuccessResponse{Data: data},
	)
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated,
		&SuccessResponse{Data: data},
	)
}

func mapErrorCodeToHTTPStatus(code apperror.ErrorCode) int {
	switch code {
	case apperror.CodeInvalidInput:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperror.CodeForbidden:
		return http.StatusForbidden
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case apperror.CodeConflict, apperror.CodeEmailTaken:
		return http.StatusConflict
	case apperror.CodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func Error(c *gin.Context, err error) {
	var appErr *apperror.AppError

	if errors.As(err, &appErr) {
		status := mapErrorCodeToHTTPStatus(appErr.Code)
		c.JSON(status, ErrorResponse{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		})
		return
	}

	// Falls back to generic 500 error if it's not our custom AppError
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    string(apperror.CodeInternal),
		Message: "Server has encountered an error, please try again later",
	})
}
