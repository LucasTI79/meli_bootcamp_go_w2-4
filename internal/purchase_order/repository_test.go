package purchaseorder_test

import (
	"context"
	"testing"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	purchaseorder "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/purchase_order"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("Creates valid purchase order", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := purchaseorder.NewRepository(db)

		order := domain.PurchaseOrder{
			OrderNumber:     "321321",
			OrderDate:       time.Now().AddDate(0, 0, 1),
			TrackingCode:    "654654",
			BuyerID:         1,
			ProductRecordID: 1,
			OrderStatusID:   1,
		}

		id, _ := repo.Create(context.TODO(), order)
		var receivedNumber string

		row := db.QueryRow(`SELECT order_number FROM purchase_orders WHERE id=?;`, id)
		row.Scan(&receivedNumber)

		assert.Equal(t, order.OrderNumber, receivedNumber)
	})
	t.Run("Does not create invalid order_number", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := purchaseorder.NewRepository(db)

		order := domain.PurchaseOrder{
			OrderDate:       time.Now().AddDate(0, 0, 1),
			TrackingCode:    "654654",
			BuyerID:         1,
			ProductRecordID: 1,
			OrderStatusID:   1,
		}

		_, err := repo.Create(context.TODO(), order)
		assert.Error(t, err)
	})
}
