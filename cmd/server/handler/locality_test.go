package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/localities"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/optional"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const LOCALITY_URL = "/localities"
const SELLER_REPORT_URL = "/report-sellers"
const CARRIER_REPORT_URL = "/report-carriers"

func TestLocalityCreate(t *testing.T) {
	t.Run("Returns 201 if locality is created successfully", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		dto := localities.CreateDTO{
			Name:     "Melicidade",
			Province: "SP",
			Country:  "BR",
		}
		expected := domain.Locality{
			ID:       1,
			Name:     "Melicidade",
			Province: "SP",
			Country:  "BR",
		}
		svc.On("Create", mock.Anything, mock.Anything).Return(expected, nil)

		req, res := testutil.MakeRequest(http.MethodPost, LOCALITY_URL, dto)
		server.ServeHTTP(res, req)

		var response testutil.SuccessResponse[domain.Locality]
		json.Unmarshal(res.Body.Bytes(), &response)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, expected, response.Data)
	})
	t.Run("Returns 409 if locality already exists", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		dto := localities.CreateDTO{
			Name:     "Melicidade",
			Province: "SP",
			Country:  "BR",
		}
		svc.On("Create", mock.Anything, mock.Anything).Return(domain.Locality{}, &localities.ErrInvalidLocality{})

		req, res := testutil.MakeRequest(http.MethodPost, LOCALITY_URL, dto)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
	})
}

func TestLocalitySellerReport(t *testing.T) {
	t.Run("Returns all localities if id is omitted", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		counts := getSellerCounts()
		expected := handler.MapSellerReportToDTO(counts)
		svc.On("CountSellers", mock.Anything, mock.Anything).Return(counts, nil)

		req, res := testutil.MakeRequest(http.MethodGet, LOCALITY_URL+SELLER_REPORT_URL, nil)
		server.ServeHTTP(res, req)

		var response testutil.SuccessResponse[[]handler.SellerReportEntry]
		json.Unmarshal(res.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.ElementsMatch(t, expected, response.Data)
	})
	t.Run("Returns single locality if id is given", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		counts := getSellerCounts()[:1]
		expected := handler.MapSellerReportToDTO(counts)
		id := *optional.FromVal(counts[0].ID)
		svc.On("CountSellers", mock.Anything, id).Return(counts, nil)

		url := fmt.Sprintf("%s/%d", LOCALITY_URL+SELLER_REPORT_URL, id.Val)
		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		var response testutil.SuccessResponse[[]handler.SellerReportEntry]
		json.Unmarshal(res.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.ElementsMatch(t, expected, response.Data)
	})
	t.Run("Returns 404 if id is not found", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		id := *optional.FromVal(42)
		svc.On("CountSellers", mock.Anything, id).Return([]localities.CountByLocality{}, localities.NewErrNotFound(id.Val))

		url := fmt.Sprintf("%s/%d", LOCALITY_URL+SELLER_REPORT_URL, id.Val)
		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
	t.Run("Returns 201 if no id is given and there is no data", func(t *testing.T) {
		t.Run("Returns 404 if id is not found", func(t *testing.T) {
			svc := LocalityServiceMock{}
			h := handler.NewLocality(&svc)
			server := getLocalityServer(h)

			id := optional.Opt[int]{}
			svc.On("CountSellers", mock.Anything, id).Return([]localities.CountByLocality{}, nil)

			req, res := testutil.MakeRequest(http.MethodGet, LOCALITY_URL+SELLER_REPORT_URL, nil)
			server.ServeHTTP(res, req)

			var response testutil.SuccessResponse[[]handler.SellerReportEntry]
			json.Unmarshal(res.Body.Bytes(), &response)

			assert.Equal(t, http.StatusNoContent, res.Code)
			assert.Len(t, response.Data, 0)
		})
	})
	t.Run("Returns 590 if repository fails", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		svc.On("CountSellers", mock.Anything, mock.Anything).Return([]localities.CountByLocality{}, localities.NewErrGeneric(""))

		req, res := testutil.MakeRequest(http.MethodGet, LOCALITY_URL+SELLER_REPORT_URL, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestLocalityCarrierReport(t *testing.T) {
	t.Run("Returns all localities if id is omitted", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		counts := getSellerCounts()
		expected := handler.MapCarrierReportToDTO(counts)
		svc.On("CountCarriers", mock.Anything, mock.Anything).Return(counts, nil)

		req, res := testutil.MakeRequest(http.MethodGet, LOCALITY_URL+CARRIER_REPORT_URL, nil)
		server.ServeHTTP(res, req)

		var response testutil.SuccessResponse[[]handler.CarrierReportEntry]
		json.Unmarshal(res.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.ElementsMatch(t, expected, response.Data)
	})
	t.Run("Returns single locality if id is given", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		counts := getSellerCounts()[:1]
		expected := handler.MapCarrierReportToDTO(counts)

		id := *optional.FromVal(counts[0].ID)
		svc.On("CountCarriers", mock.Anything, id).Return(counts, nil)

		url := fmt.Sprintf("%s/%d", LOCALITY_URL+CARRIER_REPORT_URL, id.Val)
		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		var response testutil.SuccessResponse[[]handler.CarrierReportEntry]
		json.Unmarshal(res.Body.Bytes(), &response)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.ElementsMatch(t, expected, response.Data)
	})
	t.Run("Returns 404 if id is not found", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		id := *optional.FromVal(42)
		svc.On("CountCarriers", mock.Anything, id).Return([]localities.CountByLocality{}, localities.NewErrNotFound(id.Val))

		url := fmt.Sprintf("%s/%d", LOCALITY_URL+CARRIER_REPORT_URL, id.Val)
		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
	t.Run("Returns 201 if no id is given and there is no data", func(t *testing.T) {
		t.Run("Returns 404 if id is not found", func(t *testing.T) {
			svc := LocalityServiceMock{}
			h := handler.NewLocality(&svc)
			server := getLocalityServer(h)

			id := optional.Opt[int]{}
			svc.On("CountCarriers", mock.Anything, id).Return([]localities.CountByLocality{}, nil)

			req, res := testutil.MakeRequest(http.MethodGet, LOCALITY_URL+CARRIER_REPORT_URL, nil)
			server.ServeHTTP(res, req)

			var response testutil.SuccessResponse[[]handler.CarrierReportEntry]
			json.Unmarshal(res.Body.Bytes(), &response)

			assert.Equal(t, http.StatusNoContent, res.Code)
			assert.Len(t, response.Data, 0)
		})
	})
	t.Run("Returns 590 if repository fails", func(t *testing.T) {
		svc := LocalityServiceMock{}
		h := handler.NewLocality(&svc)
		server := getLocalityServer(h)

		svc.On("CountCarriers", mock.Anything, mock.Anything).Return([]localities.CountByLocality{}, localities.NewErrGeneric(""))

		req, res := testutil.MakeRequest(http.MethodGet, LOCALITY_URL+CARRIER_REPORT_URL, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func getSellerCounts() []localities.CountByLocality {
	return []localities.CountByLocality{
		{
			ID:    1,
			Name:  "Melicidade",
			Count: 2,
		},
		{
			ID:    2,
			Name:  "Tesla",
			Count: 1,
		},
	}
}

func getLocalityServer(h *handler.Locality) *gin.Engine {
	s := testutil.CreateServer()
	rg := s.Group(LOCALITY_URL)
	{
		rg.POST("", middleware.Body[localities.CreateDTO](), h.Create())
		rg.GET(SELLER_REPORT_URL, h.SellerReport())
		rg.GET(SELLER_REPORT_URL+"/:id", middleware.IntPathParam(), h.SellerReport())
		rg.GET(CARRIER_REPORT_URL, h.CarrierReport())
		rg.GET(CARRIER_REPORT_URL+"/:id", middleware.IntPathParam(), h.CarrierReport())
	}
	return s
}

type LocalityServiceMock struct {
	mock.Mock
}

func (s *LocalityServiceMock) Create(c context.Context, loc localities.CreateDTO) (domain.Locality, error) {
	args := s.Called(c, loc)
	return args.Get(0).(domain.Locality), args.Error(1)
}

func (s *LocalityServiceMock) CountSellers(c context.Context, id optional.Opt[int]) ([]localities.CountByLocality, error) {
	args := s.Called(c, id)
	return args.Get(0).([]localities.CountByLocality), args.Error(1)
}

func (s *LocalityServiceMock) CountCarriers(c context.Context, id optional.Opt[int]) ([]localities.CountByLocality, error) {
	args := s.Called(c, id)
	return args.Get(0).([]localities.CountByLocality), args.Error(1)
}
