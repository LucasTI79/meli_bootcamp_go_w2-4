package web

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIntParam(c *gin.Context, param string) (val int, err error) {
	val64, err := strconv.ParseInt(c.Param(param), 10, 0)
	if err != nil {
		return
	}
	val = int(val64)
	return
}
