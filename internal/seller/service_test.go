package seller_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller/stubs"
	"github.com/stretchr/testify/assert"
)

func TestCreateSeller(t *testing.T) {
	t.Run("Create valid seller", func(t *testing.T) {
		stub := stubs.RepositoryStub{}
		svc := seller.NewService(&stub)

		seller := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}

		received, err := svc.Save(context.TODO(), seller)

		assert.NoError(t, err)
		assert.Equal(t, seller, received)
	})
}
