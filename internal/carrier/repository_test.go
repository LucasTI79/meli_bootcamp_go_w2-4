package carrier_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/carrier"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryCreate(t *testing.T) {
	t.Run("Creates valid carrier", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := carrier.NewRepository(db)

		c := domain.Carrier{
			CID:         873456,
			CompanyName: "meli",
			Address:     "osasco",
			Telephone:   "99999",
			LocalityID:  1,
		}

		id, _ := repo.Create(context.TODO(), c)
		var receivedCid int

		row := db.QueryRow(`SELECT cid FROM carriers WHERE id = ?;`, id)
		row.Scan(&receivedCid)

		assert.Equal(t, c.CID, receivedCid)
	})
	t.Run("Does not create invalid carrier", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := carrier.NewRepository(db)

		c := domain.Carrier{
			CID:         873456,
			CompanyName: "meli",
			Address:     "osasco",
			Telephone:   "99999",
			LocalityID:  1,
		}

		_, err := repo.Create(context.TODO(), c)
		assert.NoError(t, err)

		_, err = repo.Create(context.TODO(), c)
		assert.Error(t, err)
	})
}
