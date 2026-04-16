package register

import (
	"time"

	"github.com/gin-gonic/gin"
	coreHttp "github.com/hogiabao7725/go-ticket-engine/internal/core/delivery/http"
	authHttp "github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/delivery/http"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,not_blank"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type HTTPHandler struct {
	registerHandler *Handler
}

func NewHTTPHandler(registerHandler *Handler) *HTTPHandler {
	return &HTTPHandler{registerHandler: registerHandler}
}

func (h *HTTPHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if !coreHttp.BindJSON(c, &req) {
		return
	}

	cmd := Command{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.registerHandler.Execute(c.Request.Context(), cmd)
	if err != nil {
		status, msg := authHttp.MapDomainErrorToHTTP(err)
		coreHttp.Error(c, status, msg)
		return
	}	

	resp := RegisterResponse{
		ID:        user.ID(),
		Name:      user.Name().String(),
		Email:     user.Email().String(),
		Role:      user.Role().String(),
		CreatedAt: user.CreatedAt(),
	}

	coreHttp.Created(c, resp)
}
