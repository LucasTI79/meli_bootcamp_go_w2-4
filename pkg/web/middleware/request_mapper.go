package middleware

import (
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/gin-gonic/gin"
)

const CONTEXT_BODY_VAR_NAME = "body"

func JSONMapper[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBind(&req); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			c.Abort()
			return
		}
		c.Set(CONTEXT_BODY_VAR_NAME, req)
		c.Next()
	}
}

func ParsedRequest[T any](c *gin.Context) T {
	return c.MustGet(CONTEXT_BODY_VAR_NAME).(T)
}
