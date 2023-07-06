package carrier

import (
	"context"
	"database/sql"
	"strings"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, p domain.Carrier) (int, error)
	Exists(ctx context.Context, cid int) bool
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Exists(ctx context.Context, cid int) bool {
	query := "SELECT cid FROM carriers WHERE cid=?;"
	row := r.db.QueryRow(query, cid)
	err := row.Scan(&cid)
	return err == nil
}

func (r *repository) Create(ctx context.Context, i domain.Carrier) (int, error) {
	query := "INSERT INTO carriers(cid,company_name,address,telephone,locality_id) VALUES (?,?,?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(i.CID, i.CompanyName, i.Address, i.Telephone, i.LocalityID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1452") {
			println(err.Error())
			return 0, ErrLocalityIDNotFound
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
