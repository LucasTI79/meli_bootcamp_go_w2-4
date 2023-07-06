package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
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

		req, res := testutil.MakeRequest(http.MethodPost, CARRIER_URL, body)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[domain.Carrier]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, expected, received.Data)
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

// func getTestCarrierDTO() carrier.CarrierDTO {
// 	return carrier.CarrierDTO{
// 		CID:         10,
// 		CompanyName: "mercado livre",
// 		Address:     "osasco",
// 		Telephone:   "12345689",
// 		LocalityID:  5,
// 	}
// }

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

type CarrierServiceMock struct {
	mock.Mock
}

func (s *CarrierServiceMock) Create(c context.Context, carrier carrier.CarrierDTO) (domain.Carrier, error) {
	args := s.Called(c, carrier)
	return args.Get(0).(domain.Carrier), args.Error(1)
}
