package localities_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/localities"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryCreate(t *testing.T) {
	t.Run("Creates valid locality", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := localities.NewRepository(db)

		loc := domain.Locality{
			Name:     "Melicidade-2",
			Province: "SP",
			Country:  "BR",
		}

		id, _ := repo.Save(context.TODO(), loc)
		var receivedName string

		row := db.QueryRow(`SELECT locality_name FROM localities WHERE id = ?;`, id)
		row.Scan(&receivedName)

		assert.Equal(t, loc.Name, receivedName)
	})
	t.Run("Does not create invalid locality", func(t *testing.T) {
		db := testutil.InitDatabase(t)
		defer db.Close()

		repo := localities.NewRepository(db)

		loc := domain.Locality{
			Name:     "Melicidade-2",
			Province: "SP",
			Country:  "BR",
		}

		_, err := repo.Save(context.TODO(), loc)
		assert.NoError(t, err)

		_, err = repo.Save(context.TODO(), loc)
		assert.Error(t, err)
	})
}
