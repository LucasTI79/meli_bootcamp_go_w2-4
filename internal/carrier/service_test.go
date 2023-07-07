package carrier_test

import (
	"context"
	"errors"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/carrier"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var carrierID = 1

func TestCreate(t *testing.T) {
	t.Run("Create a carrier successfully", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := carrier.NewService(&repositoryMock)

		body := getTestCarrierDTO()

		expected := getTestCarrier()

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false)
		repositoryMock.On("Create", mock.Anything, mock.Anything).Return(carrierID, nil)

		result, err := svc.Create(context.TODO(), body)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.EqualValues(t, carrierID, result.ID)
	})
	t.Run("Does not create any carrier and returns error: cid alredy exists", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := carrier.NewService(&repositoryMock)

		body := getTestCarrierDTO()

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(true)

		_, err := svc.Create(context.TODO(), body)

		assert.ErrorIs(t, err, carrier.ErrAlreadyExists)
		repositoryMock.AssertNumberOfCalls(t, "Create", 0)
	})
	t.Run("Does not create any carrier and returns error: locality_id not found", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := carrier.NewService(&repositoryMock)

		body := getTestCarrierDTO()

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false)
		repositoryMock.On("Create", mock.Anything, mock.Anything).Return(0, carrier.ErrLocalityIDNotFound)

		_, err := svc.Create(context.TODO(), body)

		assert.ErrorIs(t, err, carrier.ErrLocalityIDNotFound)
	})
	t.Run("Does not create any carrier and returns error: internal server error", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := carrier.NewService(&repositoryMock)

		body := getTestCarrierDTO()

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false)
		repositoryMock.On("Create", mock.Anything, mock.Anything).Return(0, errors.New(""))

		_, err := svc.Create(context.TODO(), body)

		assert.ErrorIs(t, err, carrier.ErrInternalServerError)
	})
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

func getTestCarrierDTO() carrier.CarrierDTO {
	return carrier.CarrierDTO{
		CID:         10,
		CompanyName: "mercado livre",
		Address:     "osasco",
		Telephone:   "12345689",
		LocalityID:  5,
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

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Create(ctx context.Context, i domain.Carrier) (int, error) {
	args := r.Called(ctx, i)
	return args.Get(0).(int), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, cid int) bool {
	args := r.Called(ctx, cid)
	return args.Get(0).(bool)
}
