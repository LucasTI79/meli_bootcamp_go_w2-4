package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var SECTIONS_URL = "/sections/"
var sectionID = 1

func TestSectionRead(t *testing.T) {
	t.Run("Return all sections", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		expected := getTestSections()
		sectionService.On("GetAll", mock.Anything).Return(expected, nil)

		req, res := testutil.MakeRequest(http.MethodGet, SECTIONS_URL, "")
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[[]domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.ElementsMatch(t, expected, received.Data)
	})
	t.Run("Return no section and error: no content", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("GetAll", mock.Anything).Return(make([]domain.Section, 0), nil)

		req, res := testutil.MakeRequest(http.MethodGet, SECTIONS_URL, "")
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[[]domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusNoContent, res.Code)
		assert.Len(t, received.Data, 0)
	})
	t.Run("Return no section and error: internal server error", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("GetAll", mock.Anything).Return(make([]domain.Section, 0), errors.New(""))

		req, res := testutil.MakeRequest(http.MethodGet, SECTIONS_URL, "")
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[[]domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Len(t, received.Data, 0)
	})
	t.Run("Return a section by ID", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionMock := getTestSections()[0]
		sectionService.On("Get", mock.Anything, sectionID).Return(sectionMock, nil)

		url := fmt.Sprintf("%s%d", SECTIONS_URL, sectionID)
		req, res := testutil.MakeRequest(http.MethodGet, url, "")
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, sectionMock, received.Data)
	})
	t.Run("Return no section and error: not found", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("Get", mock.Anything, sectionID).Return(domain.Section{}, errors.New(""))

		url := fmt.Sprintf("%s%d", SECTIONS_URL, sectionID)
		req, res := testutil.MakeRequest(http.MethodGet, url, "")
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
}

func getSectionServer(h *handler.Section) *gin.Engine {
	server := testutil.CreateServer()

	sectionRG := server.Group(SECTIONS_URL)
	{
		sectionRG.POST("", middleware.Body[handler.CreateRequest](), h.Create())
		sectionRG.GET("", h.GetAll())
		sectionRG.GET(":id", middleware.IntPathParam(), h.Get())
		sectionRG.PATCH(":id", middleware.IntPathParam(), middleware.Body[handler.UpdateRequest](), h.Update())
		sectionRG.DELETE(":id", middleware.IntPathParam(), h.Delete())
	}

	return server
}

func getTestSections() []domain.Section {
	return []domain.Section{
		{
			ID:                 1,
			SectionNumber:      123,
			CurrentTemperature: 10,
			MinimumTemperature: 5,
			CurrentCapacity:    15,
			MinimumCapacity:    10,
			MaximumCapacity:    20,
			WarehouseID:        321,
			ProductTypeID:      2,
		},
	}
}

func getTestCreateSections() section.CreateSection {
	return section.CreateSection{
		SectionNumber:      123,
		CurrentTemperature: 10,
		MinimumTemperature: 5,
		CurrentCapacity:    15,
		MinimumCapacity:    10,
		MaximumCapacity:    20,
		WarehouseID:        321,
		ProductTypeID:      2,
	}
}

type SectionServiceMock struct {
	mock.Mock
}

func (s *SectionServiceMock) Save(ctx context.Context, section section.CreateSection) (domain.Section, error) {
	args := s.Called(ctx, section)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (s *SectionServiceMock) GetAll(ctx context.Context) ([]domain.Section, error) {
	args := s.Called(ctx)
	return args.Get(0).([]domain.Section), args.Error(1)
}

func (s *SectionServiceMock) Get(ctx context.Context, id int) (domain.Section, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (s *SectionServiceMock) Update(ctx context.Context, dto section.UpdateSection, id int) (domain.Section, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (s *SectionServiceMock) Delete(ctx context.Context, id int) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}
