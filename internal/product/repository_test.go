package product_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	product "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestRepoCreate(t *testing.T) {
	t.Run("Creates valid product", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)
		p := getTestProduct()

		id, err := repo.Save(context.TODO(), p)

		assert.NoError(t, err)

		var receivedCode string
		row := db.QueryRow(`SELECT product_code FROM products WHERE id = ?;`, id)
		row.Scan(&receivedCode)

		assert.Equal(t, p.ProductCode, receivedCode)
	})
}

func TestRepoExists(t *testing.T) {
	t.Run("Checks that product exists", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)
		p := getTestProduct()

		repo.Save(context.TODO(), p)

		assert.True(t, repo.Exists(context.TODO(), p.ProductCode))
	})
}

func TestRepoGet(t *testing.T) {
	t.Run("Gets product correctly", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)
		expected := getTestProduct()

		id, _ := repo.Save(context.TODO(), expected)

		received, err := repo.Get(context.TODO(), id)

		assert.NoError(t, err)
		assert.Equal(t, expected.ProductCode, received.ProductCode)
	})
	t.Run("Does not get inexistent product", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)

		_, err := repo.Get(context.TODO(), 999)

		assert.Error(t, err)
	})
}

func TestRepoGetAll(t *testing.T) {
	t.Run("Gets all products", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)
		expected := getTestProduct()

		repo.Save(context.TODO(), expected)

		received, err := repo.GetAll(context.TODO())

		assert.NoError(t, err)
		assert.True(t, len(received) > 0)
	})
}

func TestRepoUpdate(t *testing.T) {
	t.Run("Updates a product", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)
		p := getTestProduct()

		id, _ := repo.Save(context.TODO(), p)

		p.ID = id
		p.ProductCode = "NEW CODE"

		err := repo.Update(context.TODO(), p)
		assert.NoError(t, err)

		updated, _ := repo.Get(context.TODO(), id)
		assert.Equal(t, p.ProductCode, updated.ProductCode)
	})
}

func TestRepoDelete(t *testing.T) {
	t.Run("Deletes a product", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)
		p := getTestProduct()

		id, _ := repo.Save(context.TODO(), p)

		err := repo.Delete(context.TODO(), id)
		assert.NoError(t, err)

		_, err = repo.Get(context.TODO(), id)
		assert.Error(t, err)
	})
}

func TestRecordRepoCreate(t *testing.T) {
	t.Run("Creates valid product record", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)

		record := domain.Product_Records{
			LastUpdateDate: "2022-01-03",
			PurchasePrice:  20.20,
			SalePrice:      30.30,
			ProductID:      2,
		}

		id, _ := repo.SaveRecord(context.TODO(), record)
		var receivedNumber int

		row := db.QueryRow(`SELECT product_id FROM product_records WHERE id = ?;`, id)
		row.Scan(&receivedNumber)

		assert.Equal(t, record.ProductID, receivedNumber)
	})
	t.Run("Error on create product record repository", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)

		record := domain.Product_Records{
			LastUpdateDate: "2022-01-03",
			SalePrice:      30.30,
			ProductID:      20000,
		}

		_, err := repo.SaveRecord(context.TODO(), record)
		assert.Error(t, err)
	})
}
func TestRecordRepoGet(t *testing.T) {
	t.Run("Gets all product records", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)

		record := domain.Product_Records{
			LastUpdateDate: "2022-01-03",
			PurchasePrice:  20.20,
			SalePrice:      30.30,
			ProductID:      2,
		}

		repo.SaveRecord(context.TODO(), record)

		received, _ := repo.GetAllRecords(context.TODO())
		assert.True(t, len(received) > 0)
	})
}

func TestRecordReport(t *testing.T) {
	t.Run("Gets product reports", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := product.NewRepository(db)

		p := getTestProduct()
		id, _ := repo.Save(context.TODO(), p)
		record := domain.Product_Records{
			LastUpdateDate: "2022-01-03",
			PurchasePrice:  20.20,
			SalePrice:      30.30,
			ProductID:      id,
		}

		repo.SaveRecord(context.TODO(), record)

		_, err := repo.GetRecordsbyProd(context.TODO(), id)
		assert.NoError(t, err)
	})
}

func getTestProduct() domain.Product {
	return domain.Product{
		Description:    "abc",
		ExpirationRate: 1,
		FreezingRate:   2,
		Height:         3,
		Length:         4,
		Netweight:      5,
		ProductCode:    "PRODUCT-1",
		RecomFreezTemp: 6,
		Width:          7,
		ProductTypeID:  1,
		SellerID:       1,
	}
}
