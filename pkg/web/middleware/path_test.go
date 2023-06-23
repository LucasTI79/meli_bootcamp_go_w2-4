package middleware_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntPathParamValidator(t *testing.T) {
	t.Run("Should call handler if path param is valid", func(t *testing.T) {
		server := testutil.CreateServer()

		intParam := middleware.IntPathParam()
		handler := func(ctx *gin.Context) { web.Success(ctx, 200, nil) }
		server.GET("/:code", intParam, handler)

		code := 42
		url := fmt.Sprintf("/%d", code)
		req, res := testutil.MakeRequest(http.MethodGet, url, "")
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Should raise status 400 if path param is not an int", func(t *testing.T) {
		server := testutil.CreateServer()

		intParam := middleware.IntPathParam()
		handler := func(ctx *gin.Context) { web.Success(ctx, 200, nil) }
		server.GET("/:code", intParam, handler)

		code := "code"
		url := fmt.Sprintf("/%s", code)
		req, res := testutil.MakeRequest(http.MethodGet, url, "")
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}
