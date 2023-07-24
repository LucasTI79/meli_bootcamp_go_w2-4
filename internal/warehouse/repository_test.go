package warehouse_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/warehouse"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRepoCreate(t *testing.T) {
	t.Run("Creates valid warehouse", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := warehouse.NewRepository(db)
		wh := getTestWarehouse()

		id, err := repo.Save(context.TODO(), wh)

		assert.NoError(t, err)

		var receivedCode string
		row := db.QueryRow(`SELECT warehouse_code FROM warehouses WHERE id = ?;`, id)
		row.Scan(&receivedCode)

		assert.Equal(t, wh.WarehouseCode, receivedCode)
	})
}

func TestRepoExists(t *testing.T) {
	t.Run("Finds existing warehouse", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := warehouse.NewRepository(db)
		expected := getTestWarehouse()

		repo.Save(context.TODO(), expected)

		exists := repo.Exists(context.TODO(), expected.WarehouseCode)

		assert.True(t, exists)
	})
}

func TestRepoGet(t *testing.T) {
	t.Run("Gets warehouse", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := warehouse.NewRepository(db)
		expected := getTestWarehouse()

		id, _ := repo.Save(context.TODO(), expected)

		received, err := repo.Get(context.TODO(), id)
		assert.NoError(t, err)

		assert.Equal(t, expected.WarehouseCode, received.WarehouseCode)
	})
	t.Run("Doesn't get nonexistent warehouse", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := warehouse.NewRepository(db)

		_, err := repo.Get(context.TODO(), 9999)
		assert.Error(t, err)
	})
}

func TestRepoGetAll(t *testing.T) {
	t.Run("Gets all warehouses", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := warehouse.NewRepository(db)
		expected := getTestWarehouse()

		repo.Save(context.TODO(), expected)

		received, err := repo.GetAll(context.TODO())
		assert.NoError(t, err)

		assert.True(t, len(received) > 0)
	})
}

func TestRepoUpdate(t *testing.T) {
	t.Run("Updates a warehouse", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := warehouse.NewRepository(db)
		wh := getTestWarehouse()

		id, _ := repo.Save(context.TODO(), wh)

		wh.ID = id
		wh.WarehouseCode = "NEW CODE"

		err := repo.Update(context.TODO(), wh)
		assert.NoError(t, err)

		received, _ := repo.Get(context.TODO(), id)
		assert.Equal(t, wh.WarehouseCode, received.WarehouseCode)
	})
}

func TestRepoDelete(t *testing.T) {
	t.Run("Deletes a warehouse", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := warehouse.NewRepository(db)
		wh := getTestWarehouse()

		id, _ := repo.Save(context.TODO(), wh)

		err := repo.Delete(context.TODO(), id)
		assert.NoError(t, err)

		_, err = repo.Get(context.TODO(), id)
		assert.Error(t, err)
	})
}

func getTestWarehouse() domain.Warehouse {
	return domain.Warehouse{
		Address:            "test street",
		Telephone:          "09999",
		WarehouseCode:      "test-warehouse",
		MinimumCapacity:    1,
		MinimumTemperature: 2,
		LocalityID:         1,
	}
}
