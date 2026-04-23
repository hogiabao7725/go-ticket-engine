package refresh

import (
	"github.com/gin-gonic/gin"

	coreHttp "github.com/hogiabao7725/gin-auth-playground/internal/core/delivery/http"
	authHttp "github.com/hogiabao7725/gin-auth-playground/internal/modules/auth/delivery/http"
)

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type HTTPHandler struct {
	refreshHandler *Handler
}

func NewHTTPHandler(refreshHandler *Handler) *HTTPHandler {
	return &HTTPHandler{refreshHandler: refreshHandler}
}

func (h *HTTPHandler) Refresh(c *gin.Context) {
	token, err := c.Cookie("refresh_token")
	if err != nil || token == "" {
		coreHttp.Error(c, 401, "refresh token is missing")
		return
	}

	cmd := Command{
		RefreshToken: token,
	}

	result, err := h.refreshHandler.Execute(c.Request.Context(), cmd)
	if err != nil {
		status, msg := authHttp.MapDomainErrorToHTTP(err)
		coreHttp.Error(c, status, msg)
		return
	}

	resp := RefreshResponse{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
	}

	maxAge := int(result.RefreshExpiresIn)
	if maxAge == 0 {
		maxAge = 7 * 24 * 60 * 60 // fallback if zero
	}
	c.SetCookie("refresh_token", result.RefreshToken, maxAge, "/", "", false, true)

	coreHttp.OK(c, resp)
}
