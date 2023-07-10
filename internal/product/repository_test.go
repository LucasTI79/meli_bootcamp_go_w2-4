package product_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	product "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

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
