package batches_test

import (
	"context"
	"testing"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/batches"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryCreate(t *testing.T) {
	t.Run("Creates valid batch", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := batches.NewRepository(db)

		batch := domain.Batches{
			BatchNumber:        321874,
			CurrentQuantity:    1,
			CurrentTemperature: 2,
			DueDate:            time.Unix(100000, 1000),
			InitialQuantity:    10,
			ManufacturingDate:  time.Unix(100000, 1000),
			ManufacturingHour:  12,
			MinimumTemperature: 0,
			ProductID:          1,
			SectionID:          1,
		}

		id, _ := repo.Save(context.TODO(), batch)

		var receivedNumber int
		row := db.QueryRow(`SELECT batch_number FROM product_batches WHERE id = ?;`, id)
		row.Scan(&receivedNumber)

		assert.Equal(t, batch.BatchNumber, receivedNumber)
	})
	t.Run("Does not create invalid batch", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := batches.NewRepository(db)

		batch := domain.Batches{
			BatchNumber:        321874,
			CurrentQuantity:    1,
			CurrentTemperature: 2,
			DueDate:            time.Unix(100000, 1000),
			InitialQuantity:    10,
			ManufacturingDate:  time.Unix(100000, 1000),
			ManufacturingHour:  12,
			MinimumTemperature: 0,
			ProductID:          1,
			SectionID:          1,
		}

		_, err := repo.Save(context.TODO(), batch)
		assert.NoError(t, err)

		_, err = repo.Save(context.TODO(), batch)
		assert.Error(t, err)
	})
}
