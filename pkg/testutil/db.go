package testutil

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	txdb.Register("txdb", "mysql", "meli_sprint_user:Meli_Sprint#123@/melisprint")
}

func InitDatabase(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("txdb", uuid.New().String())
	assert.NoError(t, err)
	return db
}
