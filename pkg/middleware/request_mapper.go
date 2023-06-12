package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CONTEXT_BODY_VAR_NAME = "body"

func JSONMapper[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Split(err.Error(), "\n")})
			return
		}
		c.Set(CONTEXT_BODY_VAR_NAME, req)
		c.Next()
	}
}

func ParsedRequest[T any](c *gin.Context) T {
	return c.MustGet(CONTEXT_BODY_VAR_NAME).(T)
}
