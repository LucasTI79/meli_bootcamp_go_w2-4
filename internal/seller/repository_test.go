package seller_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRepoCreate(t *testing.T) {
	t.Run("Creates valid seller", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := seller.NewRepository(db)
		s := getTestSeller()

		id, err := repo.Save(context.TODO(), s)

		assert.NoError(t, err)

		var cid int
		row := db.QueryRow(`SELECT cid FROM sellers WHERE id = ?;`, id)
		row.Scan(&cid)

		assert.Equal(t, s.CID, cid)
	})
}
func TestRepoExists(t *testing.T) {
	t.Run("Finds existing seller", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := seller.NewRepository(db)
		expected := getTestSeller()

		repo.Save(context.TODO(), expected)

		exists := repo.Exists(context.TODO(), expected.CID)

		assert.True(t, exists)
	})
}
func TestRepoGet(t *testing.T) {
	t.Run("Gets seller", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := seller.NewRepository(db)
		expected := getTestSeller()

		id, _ := repo.Save(context.TODO(), expected)

		received, err := repo.Get(context.TODO(), id)
		assert.NoError(t, err)

		assert.Equal(t, expected.CID, received.CID)
	})
	t.Run("Doesn't get nonexistent seller", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := seller.NewRepository(db)

		_, err := repo.Get(context.TODO(), 9999)
		assert.Error(t, err)
	})
}
func TestRepoGetAll(t *testing.T) {
	t.Run("Gets all sellers", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := seller.NewRepository(db)
		expected := getTestSeller()

		repo.Save(context.TODO(), expected)

		received, err := repo.GetAll(context.TODO())
		assert.NoError(t, err)

		assert.True(t, len(received) > 0)
	})
}
func TestRepoUpdate(t *testing.T) {
	t.Run("Updates a seller", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := seller.NewRepository(db)
		s := getTestSeller()

		id, _ := repo.Save(context.TODO(), s)

		s.ID = id
		s.CID = 999

		err := repo.Update(context.TODO(), s)
		assert.NoError(t, err)

		received, _ := repo.Get(context.TODO(), id)
		assert.Equal(t, s.CID, received.CID)
	})
}
func TestRepoDelete(t *testing.T) {
	t.Run("Deletes a seller", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := seller.NewRepository(db)
		s := getTestSeller()

		id, _ := repo.Save(context.TODO(), s)

		err := repo.Delete(context.TODO(), id)
		assert.NoError(t, err)

		_, err = repo.Get(context.TODO(), id)
		assert.Error(t, err)
	})
}

func getTestSeller() domain.Seller {
	return domain.Seller{
		CID:         1,
		CompanyName: "meli",
		Address:     "osasco",
		Telephone:   "123",
		LocalityID:  1,
	}
}
