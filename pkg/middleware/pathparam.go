package middleware

import (
	"errors"
	"net/http"
	"strconv"

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
			c.AbortWithError(http.StatusBadRequest, errors.New("path parameter should be an int"))
			return
		}

		val := int(val64)
		c.Set(p.Key, val)

		c.Next()
	}
}
