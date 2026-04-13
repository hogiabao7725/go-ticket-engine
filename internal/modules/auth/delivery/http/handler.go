package http

import (
	"github.com/gin-gonic/gin"
	corehttp "github.com/hogiabao7725/go-ticket-engine/internal/core/delivery/http"
	"github.com/hogiabao7725/go-ticket-engine/internal/modules/auth/usecase"
)

type AuthHandler struct {
	registerer Registerer
}

func NewAuthHandler(reg Registerer) *AuthHandler {
	return &AuthHandler{
		registerer: reg,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var reqHTTP registerRequestHTTP
	if !corehttp.BindJSON(c, &reqHTTP) {
		return
	}

	// HTTP DTO -> usecase DTO
	req := usecase.RegisterRequest{
		Name:     reqHTTP.Name,
		Email:    reqHTTP.Email,
		Password: reqHTTP.Password,
	}

	resp, err := h.registerer.Execute(c.Request.Context(), req)
	if err != nil {
		status, message := mapDomainErrorToHTTP(err)
		corehttp.Error(c, status, message)
		return
	}

	// usecase DTO -> HTTP DTO
	corehttp.Created(c, registerResponseHTTP{
		ID:        resp.ID,
		Name:      resp.Name,
		Email:     resp.Email,
		Role:      resp.Role,
		CreatedAt: resp.CreatedAt,
	})
}
