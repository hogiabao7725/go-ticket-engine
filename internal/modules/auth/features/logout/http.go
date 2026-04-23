package logout

import (
	"github.com/gin-gonic/gin"

	coreHttp "github.com/hogiabao7725/gin-auth-playground/internal/core/delivery/http"
	authHttp "github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/delivery/http"
)

type HTTPHandler struct {
	logoutHandler *Handler
}

func NewHTTPHandler(logoutHandler *Handler) *HTTPHandler {
	return &HTTPHandler{logoutHandler: logoutHandler}
}

func (h *HTTPHandler) Logout(c *gin.Context) {
	token, _ := c.Cookie("refresh_token")

	cmd := Command{
		RefreshToken: token,
	}

	err := h.logoutHandler.Execute(c.Request.Context(), cmd)
	if err != nil {
		status, msg := authHttp.MapDomainErrorToHTTP(err)
		coreHttp.Error(c, status, msg)
		return
	}

	// Clear the refresh token cookie
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	coreHttp.OK(c, map[string]string{"message": "logged out successfully"})
}
