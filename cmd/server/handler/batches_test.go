package handler_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/batches"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var BATCHES_URL = "/product-batches"

func TestCreateBatches(t *testing.T) {
	t.Run("should return error when have error in the server", func(t *testing.T) {
		batchesServiceMock := BatchesServiceMock{}
		h := handler.NewBatches(&batchesServiceMock)
		server := getBatchesServer(h)

		fakeStruct := handler.CreateBatchesRequest{
			BatchNumber:        1,
			CurrentQuantity:    200,
			CurrentTemperature: 20,
			DueDate:            "2020-04-04",
			InitialQuantity:    10,
			ManufacturingDate:  "2020-04-04",
			ManufacturingHour:  10,
			MinimumTemperature: 5,
			ProductID:          2,
			SectionID:          1,
		}

		expectedDate, _ := handler.ConvertDate(fakeStruct)

		batchesServiceMock.On("Create", mock.Anything, expectedDate).Return(domain.Batches{}, errors.New("error"))
		request, response := testutil.MakeRequest(http.MethodPost, BATCHES_URL, fakeStruct)

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
	t.Run("should return error when convertdate not convert date", func(t *testing.T) {
		batchesServiceMock := BatchesServiceMock{}
		handler := handler.NewBatches(&batchesServiceMock)
		server := getBatchesServer(handler)

		fakeStruct := batches.CreateBatches{
			BatchNumber:        1,
			CurrentQuantity:    200,
			CurrentTemperature: 20,
			DueDate:            time.Date(2020, time.April, 4, 10, 0, 0, 0, time.UTC),
			InitialQuantity:    10,
			ManufacturingDate:  time.Date(2020, time.April, 4, 10, 0, 0, 0, time.UTC),
			ManufacturingHour:  10,
			MinimumTemperature: 5,
			ProductID:          2,
			SectionID:          1,
		}

		batchesServiceMock.On("Create", mock.Anything, fakeStruct).Return(domain.Batches{}, errors.New("error"))
		request, response := testutil.MakeRequest(http.MethodPost, BATCHES_URL, fakeStruct)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusUnprocessableEntity)
	})

	t.Run("should return error when convertdate convert date, but batchesNumber alread exist ", func(t *testing.T) {
		batchesServiceMock := BatchesServiceMock{}
		h := handler.NewBatches(&batchesServiceMock)
		server := getBatchesServer(h)

		fakeStruct := handler.CreateBatchesRequest{
			BatchNumber:        1,
			CurrentQuantity:    200,
			CurrentTemperature: 20,
			DueDate:            "2020-04-04",
			InitialQuantity:    10,
			ManufacturingDate:  "2020-04-04",
			ManufacturingHour:  10,
			MinimumTemperature: 5,
			ProductID:          2,
			SectionID:          1,
		}

		expectedDate, _ := handler.ConvertDate(fakeStruct)

		batchesServiceMock.On("Create", mock.Anything, expectedDate).Return(domain.Batches{}, batches.ErrInvalidBatchNumber)
		request, response := testutil.MakeRequest(http.MethodPost, BATCHES_URL, fakeStruct)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusConflict)
	})

	t.Run("should return 201 for sucess", func(t *testing.T) {
		batchesServiceMock := BatchesServiceMock{}
		h := handler.NewBatches(&batchesServiceMock)
		server := getBatchesServer(h)

		fakeStruct := handler.CreateBatchesRequest{
			BatchNumber:        1,
			CurrentQuantity:    200,
			CurrentTemperature: 20,
			DueDate:            "2020-04-04",
			InitialQuantity:    10,
			ManufacturingDate:  "2020-04-04",
			ManufacturingHour:  10,
			MinimumTemperature: 5,
			ProductID:          2,
			SectionID:          1,
		}

		expectedDate, _ := handler.ConvertDate(fakeStruct)

		batchesServiceMock.On("Create", mock.Anything, expectedDate).Return(domain.Batches{}, nil)
		request, response := testutil.MakeRequest(http.MethodPost, BATCHES_URL, fakeStruct)

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusCreated)
	})

	t.Run("should return error when convertdate not convert date (ManufacturingDate)", func(t *testing.T) {
		fakeStruct := handler.CreateBatchesRequest{
			BatchNumber:        1,
			CurrentQuantity:    200,
			CurrentTemperature: 20,
			DueDate:            "2020-04-04",
			InitialQuantity:    10,
			ManufacturingDate:  "abcde",
			ManufacturingHour:  10,
			MinimumTemperature: 5,
			ProductID:          2,
			SectionID:          1,
		}

		_, err := handler.ConvertDate(fakeStruct)
		assert.Error(t, err)
	})

}

func getBatchesServer(h *handler.Batches) *gin.Engine {
	server := testutil.CreateServer()

	server.POST(BATCHES_URL, middleware.Body[handler.CreateBatchesRequest](), h.Create())

	return server
}

type BatchesServiceMock struct {
	mock.Mock
}

func (m *BatchesServiceMock) Create(ctx context.Context, b batches.CreateBatches) (domain.Batches, error) {
	args := m.Called(ctx, b)
	return args.Get(0).(domain.Batches), args.Error(1)
}
