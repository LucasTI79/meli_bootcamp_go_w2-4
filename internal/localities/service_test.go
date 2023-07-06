package localities_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/localities"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	t.Run("Creates valid locality", func(t *testing.T) {
		repo := RepositoryMock{}
		svc := localities.NewService(&repo)

		dto := localities.CreateDTO{
			Name:     "Melicidade",
			Province: "SP",
			Country:  "Brasil",
		}
		expected := domain.Locality{
			ID:       1,
			Name:     "Melicidade",
			Province: "SP",
			Country:  "Brasil",
		}
		repo.On("Save", mock.Anything, mock.Anything).Return(1, nil)

		received, err := svc.Create(context.TODO(), dto)

		assert.NoError(t, err)
		assert.Equal(t, expected, received)
	})
	t.Run("Doesn't create duplicate locality", func(t *testing.T) {
		repo := RepositoryMock{}
		svc := localities.NewService(&repo)

		dto := localities.CreateDTO{
			Name:     "Melicidade",
			Province: "SP",
			Country:  "Brasil",
		}
		loc := domain.Locality{
			Name:     "Melicidade",
			Province: "SP",
			Country:  "Brasil",
		}
		var expectedErr *localities.ErrInvalidLocality
		repo.On("Save", mock.Anything, mock.Anything).Return(0, localities.NewErrInvalidLocality(loc))

		_, err := svc.Create(context.TODO(), dto)

		assert.ErrorAs(t, err, &expectedErr)
	})
}

func TestReport(t *testing.T) {
	t.Run("Returns all localities when id is omitted", func(t *testing.T) {
		repo := RepositoryMock{}
		svc := localities.NewService(&repo)

		noID := optional.Opt[int]{}

		expected := []localities.SellersByLocality{
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
		counts := []localities.SellerCount{
			{
				LocalityID:  1,
				SellerCount: 2,
			},
			{
				LocalityID:  2,
				SellerCount: 1,
			},
		}
		locs := getLocalities()

		repo.On("GetAll", mock.Anything).Return(locs, nil)
		repo.On("CountSellersByLocalities", mock.Anything, mock.Anything).Return(counts, nil)

		received, err := svc.CountSellers(context.TODO(), noID)

		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, received)
	})
	t.Run("Returns specific locality if id is provided", func(t *testing.T) {
		repo := RepositoryMock{}
		svc := localities.NewService(&repo)

		id := *optional.FromVal(2)

		expected := []localities.SellersByLocality{
			{
				ID:    2,
				Name:  "Tesla",
				Count: 1,
			},
		}
		counts := []localities.SellerCount{
			{
				LocalityID:  2,
				SellerCount: 1,
			},
		}
		locs := getLocalities()

		repo.On("GetAll", mock.Anything).Return(locs, nil)
		repo.On("CountSellersByLocalities", mock.Anything, []int{id.Val}).Return(counts, nil)

		received, err := svc.CountSellers(context.TODO(), id)

		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, received)
	})
	t.Run("Returns NotFound if id doesn't exist", func(t *testing.T) {
		repo := RepositoryMock{}
		svc := localities.NewService(&repo)

		id := *optional.FromVal(7)

		counts := []localities.SellerCount{}
		locs := getLocalities()
		var expectedErr *localities.ErrNotFound

		repo.On("GetAll", mock.Anything).Return(locs, nil)
		repo.On("CountSellersByLocalities", mock.Anything, []int{id.Val}).Return(counts, nil)

		_, err := svc.CountSellers(context.TODO(), id)

		assert.ErrorAs(t, err, &expectedErr)
	})
	t.Run("Doesn't return NotFound if id is omitted but no data exists", func(t *testing.T) {
		repo := RepositoryMock{}
		svc := localities.NewService(&repo)

		noID := optional.Opt[int]{}

		counts := []localities.SellerCount{}
		locs := getLocalities()

		repo.On("GetAll", mock.Anything).Return(locs, nil)
		repo.On("CountSellersByLocalities", mock.Anything, mock.Anything).Return(counts, nil)

		received, err := svc.CountSellers(context.TODO(), noID)

		assert.NoError(t, err)
		assert.Len(t, received, 0)
	})
}

func getLocalities() []domain.Locality {
	return []domain.Locality{
		{
			ID:   1,
			Name: "Melicidade",
		},
		{
			ID:   2,
			Name: "Tesla",
		},
	}
}

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Save(c context.Context, loc domain.Locality) (int, error) {
	args := r.Called(c, loc)
	return args.Get(0).(int), args.Error(1)
}

func (r *RepositoryMock) GetAll(c context.Context) ([]domain.Locality, error) {
	args := r.Called(c)
	return args.Get(0).([]domain.Locality), args.Error(1)
}

func (r *RepositoryMock) CountSellersByLocalities(c context.Context, ids []int) ([]localities.SellerCount, error) {
	args := r.Called(c, ids)
	return args.Get(0).([]localities.SellerCount), args.Error(1)
}
