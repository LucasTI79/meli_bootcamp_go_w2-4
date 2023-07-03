package middleware

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/gin-gonic/gin"
)

// Converts a path parameter to an integer.
//
//	Note:this middleware expects the endpoint
//	to have exactly one path parameter, otherwise
//	a panic occurs.
func IntPathParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Params) != 1 {
			panic("endpoint should have exactly one path parameter")
		}

		p := c.Params[0]
		val64, err := strconv.ParseInt(p.Value, 10, 0)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "path parameter %s should be an int", p.Key)
			c.Abort()
			return
		}

		val := int(val64)
		c.Set(p.Key, val)

		c.Next()
	}
}
