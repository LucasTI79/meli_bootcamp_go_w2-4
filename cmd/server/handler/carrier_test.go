package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/carrier"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var carrierID = 1
var CARRIER_URL = "/carrier/"

//var CARRIER_URL_ID = fmt.Sprintf("%s/%d", CARRIER_URL, 1)

func TestCarrierCreate(t *testing.T) {
	t.Run("Create a carrier successfully", func(t *testing.T) {
		carrierService := CarrierServiceMock{}
		h := handler.NewCarrier(&carrierService)
		server := getCarrierServer(h)

		body := getTestCarrierRequest()
		expected := getTestCarrier()

		carrierService.On("Create", mock.Anything, mock.Anything).Return(expected, nil)

		res := requestCarrierPost(body, server)

		var received testutil.SuccessResponse[domain.Carrier]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, expected, received.Data)
	})
	t.Run("Does not create any carrier and returns error: unprocessable content", func(t *testing.T) {
		carrierService := CarrierServiceMock{}
		h := handler.NewCarrier(&carrierService)
		server := getCarrierServer(h)

		body := handler.CarrierRequest{}
		res := requestCarrierPost(body, server)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
		carrierService.AssertNumberOfCalls(t, "Create", 0)
	})
	t.Run("Does not create any carrier and returns error: conflict", func(t *testing.T) {
		carrierService := CarrierServiceMock{}
		h := handler.NewCarrier(&carrierService)
		server := getCarrierServer(h)

		carrierService.On("Create", mock.Anything, mock.Anything).Return(domain.Carrier{}, carrier.ErrAlreadyExists)

		body := getTestCarrierRequest()
		res := requestCarrierPost(body, server)

		assert.Equal(t, http.StatusConflict, res.Code)
	})
	t.Run("Does not create any carrier and returns error: internal server error", func(t *testing.T) {
		carrierService := CarrierServiceMock{}
		h := handler.NewCarrier(&carrierService)
		server := getCarrierServer(h)

		carrierService.On("Create", mock.Anything, mock.Anything).Return(domain.Carrier{}, errors.New(""))

		body := getTestCarrierRequest()
		res := requestCarrierPost(body, server)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

}

func getCarrierServer(h *handler.Carrier) *gin.Engine {
	server := testutil.CreateServer()

	carrierRG := server.Group(CARRIER_URL)
	{
		carrierRG.POST("", middleware.Body[handler.CarrierRequest](), h.Create())
	}

	return server
}

func getTestCarrierRequest() handler.CarrierRequest {
	return handler.CarrierRequest{
		CID:         testutil.ToPtr(10),
		CompanyName: testutil.ToPtr("mercado livre"),
		Address:     testutil.ToPtr("osasco"),
		Telephone:   testutil.ToPtr("123456789"),
		LocalityID:  testutil.ToPtr(5),
	}
}

func getTestCarrier() domain.Carrier {
	return domain.Carrier{
		ID:          carrierID,
		CID:         10,
		CompanyName: "mercado livre",
		Address:     "osasco",
		Telephone:   "12345689",
		LocalityID:  5,
	}
}

func requestCarrierPost(body any, server *gin.Engine) *httptest.ResponseRecorder {
	req, res := testutil.MakeRequest(http.MethodPost, CARRIER_URL, body)
	server.ServeHTTP(res, req)
	return res
}

type CarrierServiceMock struct {
	mock.Mock
}

func (s *CarrierServiceMock) Create(c context.Context, carrier carrier.CarrierDTO) (domain.Carrier, error) {
	args := s.Called(c, carrier)
	return args.Get(0).(domain.Carrier), args.Error(1)
}
