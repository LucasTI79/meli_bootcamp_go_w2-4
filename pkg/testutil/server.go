package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func CreateServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func MakeRequest(method, url string, body any) (*http.Request, *httptest.ResponseRecorder) {
	marshalled, _ := json.Marshal(body)
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(marshalled)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func ToPtr[T any](val T) *T {
	return &val
}
