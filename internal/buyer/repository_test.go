package buyer_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	t.Run("Get all buyers", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		buyers, _ := repo.GetAll(context.Background())
		assert.Equal(t, 2, len(buyers))
	})
}
func TestGet(t *testing.T) {
	t.Run("Get buyer by id", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		buyer, _ := repo.Get(context.Background(), 1)
		assert.Equal(t, 1, buyer.ID)
	})
	t.Run("Error when get buyer by id", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		_, err := repo.Get(context.Background(), 100)
		assert.Error(t, err)
	})
}
func TestUpdate(t *testing.T) {
	t.Run("Update buyer", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		buyer := domain.Buyer{
			ID:           1,
			CardNumberID: "1234567890",
			FirstName:    "Juan",
			LastName:     "Perez",
		}
		err := repo.Update(context.Background(), buyer)
		assert.Nil(t, err)
		buyer, _ = repo.Get(context.Background(), 1)
		assert.Equal(t, "Juan", buyer.FirstName)
	})
}
func TestDelete(t *testing.T) {
	t.Run("error when delete buyer", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		err := repo.Delete(context.Background(), 1)
		assert.Error(t, err)
	})
	t.Run("delete buyer", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		err := repo.Delete(context.Background(), 2)
		assert.NotNil(t, err)
	})
}
func TestExists(t *testing.T) {
	t.Run("buyer exists", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		exists := repo.Exists(context.Background(), "123456789")
		assert.True(t, exists)
	})
	t.Run("buyer not exists", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		exists := repo.Exists(context.Background(), "1234567890")
		assert.False(t, exists)
	})
}
func TestSave(t *testing.T) {
	t.Run("save buyer", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		buyer := domain.Buyer{
			CardNumberID: "1234567890",
			FirstName:    "Juan",
			LastName:     "Perez",
		}
		_, err := repo.Save(context.Background(), buyer)
		assert.Nil(t, err)
	})
	t.Run("error when save buyer", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		buyer := domain.Buyer{
			CardNumberID: "123456789",
			FirstName:    "Juan",
			LastName:     "Perez",
		}
		_, err := repo.Save(context.Background(), buyer)
		assert.Error(t, err)
	})
}
func TestGetAllPurchaseOrders(t *testing.T) {
	t.Run("Get all purchase orders", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		purchaseOrders, _ := repo.GetAllPurchaseOrders(context.Background())
		assert.Equal(t, 2, len(purchaseOrders))
	})
}
func TestGetPurchaseOrderByID(t *testing.T) {
	t.Run("Get purchase order by id", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		purchaseOrder, _ := repo.GetPurchaseOrderByID(context.Background(), 1)
		assert.Equal(t, 1, purchaseOrder.ID)
	})
	t.Run("Error when get purchase order by id", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()
		repo := buyer.NewRepository(db)
		_, err := repo.GetPurchaseOrderByID(context.Background(), 100)
		assert.Error(t, err)
	})
}
