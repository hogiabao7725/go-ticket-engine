package login

import (
	"errors"

	"github.com/gin-gonic/gin"

	coreHttp "github.com/hogiabao7725/go-ticket-engine/internal/core/delivery/http"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/domain"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	ExpiresIn   int64        `json:"expires_in"`
	User        userResponse `json:"user"`
}

type HTTPHandler struct {
	loginHandler *Handler
}

func NewHTTPHandler(loginHandler *Handler) *HTTPHandler {
	return &HTTPHandler{loginHandler: loginHandler}
}

func (h *HTTPHandler) Login(c *gin.Context) {
	var req LoginRequest
	if !coreHttp.BindJSON(c, &req) {
		return
	}

	cmd := Command{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.loginHandler.Execute(c.Request.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			coreHttp.Unauthorized(c, "invalid email or password")
			return
		}

		// c.Error(err) // for logging
		coreHttp.Error(c, 500, "internal server error")
		return
	}

	resp := LoginResponse{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
		User: userResponse{
			ID:    result.User.ID(),
			Name:  result.User.Name().String(),
			Email: result.User.Email().String(),
			Role:  result.User.Role().String(),
		},
	}

	coreHttp.OK(c, resp)
}
