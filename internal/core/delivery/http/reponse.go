package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type successResponse struct {
	Data any            `json:"data"`
	Meta map[string]any `json:"meta,omitempty"`
}

type errorResponse struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, successResponse{Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, successResponse{Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, errorResponse{Message: message})
}

func ValidationError(c *gin.Context, fields map[string]string) {
	c.JSON(http.StatusBadRequest, errorResponse{
		Message: "validation failed",
		Fields:  fields,
	})
}
