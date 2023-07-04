package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

var sectionID = 1
var SECTIONS_URL = "/sections"

// Units tests
func TestSectionRead(t *testing.T) {
	t.Run("Return all sections successfully", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		expected := getTestSections()
		sectionService.On("GetAll", mock.Anything).Return(expected, nil)

		res := requestGet(server, SECTIONS_URL)

		var received testutil.SuccessResponse[[]domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.ElementsMatch(t, expected, received.Data)
	})
	t.Run("Does not get any section and returns error: no content", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("GetAll", mock.Anything).Return(make([]domain.Section, 0), nil)

		res := requestGet(server, SECTIONS_URL)

		var received testutil.SuccessResponse[[]domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusNoContent, res.Code)
		assert.Len(t, received.Data, 0)
	})
	t.Run("Does not get any section and returns error: internal server error", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("GetAll", mock.Anything).Return(make([]domain.Section, 0), errors.New(""))

		res := requestGet(server, SECTIONS_URL)

		var received testutil.SuccessResponse[[]domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Len(t, received.Data, 0)
	})
	t.Run("Return a section by ID successfully", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		expected := getTestSections()[0]
		sectionService.On("Get", mock.Anything, sectionID).Return(expected, nil)

		urlWithID := fmt.Sprintf("%s/%d", SECTIONS_URL, 1)
		res := requestGet(server, urlWithID)

		var received testutil.SuccessResponse[domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expected, received.Data)
	})
	t.Run("Does not get any section and returns error: not found", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("Get", mock.Anything, sectionID).Return(domain.Section{}, errors.New(""))

		urlWithID := fmt.Sprintf("%s/%d", SECTIONS_URL, 1)
		res := requestGet(server, urlWithID)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
}

func TestSectionCreate(t *testing.T) {
	t.Run("Create a section successfully", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		body := getTestCreateSection()
		expected := getTestSections()[0]

		sectionService.On("Create", mock.Anything, body).Return(expected, nil)

		res := requestPost(body, server, SECTIONS_URL)

		var received testutil.SuccessResponse[domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, expected, received.Data)
	})
	t.Run("Does not create any section and returns error: unprocessable content", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		body := section.CreateSection{}
		res := requestPost(body, server, SECTIONS_URL)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
		sectionService.AssertNumberOfCalls(t, "Create", 0)
	})
	t.Run("Does not create any section and returns error: conflict", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("Create", mock.Anything, mock.Anything).Return(domain.Section{}, section.ErrInvalidSectionNumber)

		body := getTestCreateSection()
		res := requestPost(body, server, SECTIONS_URL)

		assert.Equal(t, http.StatusConflict, res.Code)
	})
	t.Run("Does not create any section and returns error: internal server error", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("Create", mock.Anything, mock.Anything).Return(domain.Section{}, errors.New(""))

		body := getTestCreateSection()
		res := requestPost(body, server, SECTIONS_URL)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestSectionUpdate(t *testing.T) {
	t.Run("Update a section successfully", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		expected := getTestSections()[0]
		expected.CurrentTemperature = 11

		sectionService.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(expected, nil)

		body := getUpdateSection()
		urlWithID := fmt.Sprintf("%s/%d", SECTIONS_URL, expected.ID)
		res := requestPatch(body, server, urlWithID)

		var received testutil.SuccessResponse[domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, expected, received.Data)
	})
	t.Run("Does not update any section and returns error: not found", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(domain.Section{}, section.ErrNotFound)

		body := getUpdateSection()
		urlWithID := fmt.Sprintf("%s/%d", SECTIONS_URL, 1)
		res := requestPatch(body, server, urlWithID)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
	t.Run("Does not update any section and returns error: unprocessable content", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		body := []section.UpdateSection{}
		urlWithID := fmt.Sprintf("%s/%d", SECTIONS_URL, 1)
		res := requestPatch(body, server, urlWithID)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
		sectionService.AssertNumberOfCalls(t, "Update", 0)
	})
	t.Run("Does not update any section and returns error: conflict", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(domain.Section{}, section.ErrInvalidSectionNumber)

		body := getUpdateSection()
		urlWithID := fmt.Sprintf("%s/%d", SECTIONS_URL, 1)
		res := requestPatch(body, server, urlWithID)

		assert.Equal(t, http.StatusConflict, res.Code)
	})
	t.Run("Does not update any section and returns error: internal server error", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		body := getUpdateSection()

		sectionService.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(domain.Section{}, errors.New(""))

		urlWithID := fmt.Sprintf("%s/%d", SECTIONS_URL, 1)
		res := requestPatch(body, server, urlWithID)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func TestSectionDelete(t *testing.T) {
	t.Run("Delete a section successfully", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("Delete", mock.Anything, sectionID).Return(nil)

		urlWithID := fmt.Sprintf("%s/%d", SECTIONS_URL, 1)
		res := requestDelete(server, urlWithID)

		var received testutil.SuccessResponse[domain.Section]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Empty(t, received.Data)
		assert.Equal(t, http.StatusNoContent, res.Code)
	})
	t.Run("Does not delete any section and returns error: not found", func(t *testing.T) {
		sectionService := SectionServiceMock{}
		h := handler.NewSection(&sectionService)
		server := getSectionServer(h)

		sectionService.On("Delete", mock.Anything, mock.Anything).Return(errors.New(""))

		urlWithID := fmt.Sprintf("%s/%d", SECTIONS_URL, 1)
		res := requestDelete(server, urlWithID)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
}

// Requests

func requestGet(server *gin.Engine, url string) *httptest.ResponseRecorder {
	req, res := testutil.MakeRequest(http.MethodGet, url, "")
	server.ServeHTTP(res, req)
	return res
}

func requestPost(body section.CreateSection, server *gin.Engine, url string) *httptest.ResponseRecorder {
	req, res := testutil.MakeRequest(http.MethodPost, url, body)
	server.ServeHTTP(res, req)
	return res
}

func requestPatch(body any, server *gin.Engine, url string) *httptest.ResponseRecorder {
	req, res := testutil.MakeRequest(http.MethodPatch, url, body)
	server.ServeHTTP(res, req)
	return res
}

func requestDelete(server *gin.Engine, url string) *httptest.ResponseRecorder {
	req, res := testutil.MakeRequest(http.MethodDelete, url, "")
	server.ServeHTTP(res, req)
	return res
}

// Generate test objects

func getSectionServer(h *handler.Section) *gin.Engine {
	server := testutil.CreateServer()

	sectionRG := server.Group(SECTIONS_URL)
	{
		sectionRG.POST("", middleware.Body[section.CreateSection](), h.Create())
		sectionRG.GET("", h.GetAll())
		sectionRG.GET("/:id", middleware.IntPathParam(), h.Get())
		sectionRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
		sectionRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[section.UpdateSection](), h.Update())
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

func getTestCreateSection() section.CreateSection {
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

func getUpdateSection() section.UpdateSection {
	return section.UpdateSection{
		SectionNumber:      testutil.ToPtr(123),
		CurrentTemperature: testutil.ToPtr(11),
	}
}

// Mock service functions

type SectionServiceMock struct {
	mock.Mock
}

func (s *SectionServiceMock) Create(ctx context.Context, section section.CreateSection) (domain.Section, error) {
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
	args := s.Called(ctx, dto, id)
	return args.Get(0).(domain.Section), args.Error(1)
}

func (s *SectionServiceMock) Delete(ctx context.Context, id int) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}
