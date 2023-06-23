package middleware_test

import (
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBodyMapper(t *testing.T) {
	t.Run("Should proceed to handler if body parsing succeeds", func(t *testing.T) {
		server := testutil.CreateServer() // TODO: Remove testutil dependency

		bodyMapper := middleware.Body[PrimitiveTypesBody]()
		handler := func(ctx *gin.Context) { web.Success(ctx, 200, nil) }
		server.POST("/", bodyMapper, handler)

		body := map[string]any{
			"id":     23,
			"name":   "John Doe",
			"height": 1.91,
		}
		req, res := testutil.MakeRequest(http.MethodPost, "/", body)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Should return 422 if body has invalid type", func(t *testing.T) {
		server := testutil.CreateServer() // TODO: Remove testutil dependency

		bodyMapper := middleware.Body[PrimitiveTypesBody]()
		handler := func(ctx *gin.Context) { web.Success(ctx, 200, nil) }
		server.POST("/", bodyMapper, handler)

		body := map[string]any{
			"id":     "23",
			"name":   "John Doe",
			"height": 1.91,
		}
		req, res := testutil.MakeRequest(http.MethodPost, "/", body)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
	})
}

type PrimitiveTypesBody struct {
	ID     int     `binding:"required"`
	Name   string  `binding:"required"`
	Height float64 `binding:"required"`
}
