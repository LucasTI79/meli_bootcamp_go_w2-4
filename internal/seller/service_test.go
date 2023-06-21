package seller_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSeller(t *testing.T) {
	t.Run("Create valid seller", func(t *testing.T) {
		repositoryMock := mocks.RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		seller := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}

		repositoryMock.On("Exists", mock.Anything, 123).Return(false)
		repositoryMock.On("Save", mock.Anything, seller).Return(1, nil)

		received, err := svc.Save(context.TODO(), seller)

		assert.NoError(t, err)
		assert.Equal(t, seller, received)
	})
}
