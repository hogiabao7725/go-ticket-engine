package get_me

import (
	"github.com/gin-gonic/gin"
	coreHttp "github.com/hogiabao7725/gin-auth-playground/internal/core/delivery/http"
	"github.com/hogiabao7725/gin-auth-playground/internal/core/middleware"
	authHttp "github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/delivery/http"
)

type HTTPHandler struct {
	handler *Handler
}

func NewHTTPHandler(handler *Handler) *HTTPHandler {
	return &HTTPHandler{handler: handler}
}

func (h *HTTPHandler) GetMe(c *gin.Context) {
	userIDStr := middleware.GetUserID(c.Request.Context())
	if userIDStr == "" {
		coreHttp.Unauthorized(c, "unauthorized")
		return
	}

	cmd := Command{UserID: userIDStr}
	userDTO, err := h.handler.Execute(c.Request.Context(), cmd)
	if err != nil {
		status, msg := authHttp.MapDomainErrorToHTTP(err)
		coreHttp.Error(c, status, msg)
		return
	}

	coreHttp.OK(c, userDTO)
}
