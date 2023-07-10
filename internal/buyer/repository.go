package buyer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Repository encapsulates the storage of a buyer.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Buyer, error)
	Get(ctx context.Context, id int) (domain.Buyer, error)
	Exists(ctx context.Context, cardNumberID string) bool
	Save(ctx context.Context, b domain.Buyer) (int, error)
	Update(ctx context.Context, b domain.Buyer) error
	Delete(ctx context.Context, id int) error
	GetAllPurchaseOrders(ctx context.Context) ([]CountByBuyer, error)
	GetPurchaseOrderByID(ctx context.Context, id int) (CountByBuyer, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	query := "SELECT * FROM buyers"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var buyers []domain.Buyer

	for rows.Next() {
		b := domain.Buyer{}
		_ = rows.Scan(&b.ID, &b.CardNumberID, &b.FirstName, &b.LastName)
		buyers = append(buyers, b)
	}

	return buyers, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Buyer, error) {
	query := "SELECT * FROM buyers WHERE id = ?;"
	row := r.db.QueryRow(query, id)
	b := domain.Buyer{}
	err := row.Scan(&b.ID, &b.CardNumberID, &b.FirstName, &b.LastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Buyer{}, ErrNotFound
		}
		return domain.Buyer{}, err
	}

	return b, nil
}

func (r *repository) Exists(ctx context.Context, cardNumberID string) bool {
	query := "SELECT card_number_id FROM buyers WHERE card_number_id=?;"
	row := r.db.QueryRow(query, cardNumberID)
	err := row.Scan(&cardNumberID)
	return err == nil
}

func (r *repository) Save(ctx context.Context, b domain.Buyer) (int, error) {
	query := "INSERT INTO buyers(card_number_id,first_name,last_name) VALUES (?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&b.CardNumberID, &b.FirstName, &b.LastName)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, b domain.Buyer) error {
	query := "UPDATE buyers SET first_name=?, last_name=?  WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&b.FirstName, &b.LastName, &b.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return ErrNotFound
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM buyers WHERE id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected < 1 {
		return ErrNotFound
	}

	return nil
}

func (r *repository) GetAllPurchaseOrders(ctx context.Context) ([]CountByBuyer, error) {
	const query = `SELECT e.id, e.card_number_id, e.first_name, e.last_name, COUNT(i.id) as purchase_orders_count 
	FROM buyers e 
	LEFT JOIN purchase_orders i ON i.buyer_id = e.id 
	GROUP BY e.id;`

	rows, err := r.db.Query(query)
	if err != nil {
		return []CountByBuyer{}, ErrInternalServerError
	}
	defer rows.Close()

	var reports []CountByBuyer
	for rows.Next() {
		var e CountByBuyer
		err := rows.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.Count)
		if err != nil {
			return []CountByBuyer{}, ErrInternalServerError
		}
		reports = append(reports, e)
	}
	err = rows.Err()
	if err != nil {
		return []CountByBuyer{}, ErrInternalServerError
	}

	if len(reports) == 0 {
		return []CountByBuyer{}, ErrNotFound
	}

	return reports, nil
}

func (r *repository) GetPurchaseOrderByID(ctx context.Context, id int) (CountByBuyer, error) {
	const query = `SELECT e.id, e.card_number_id, e.first_name, e.last_name, COUNT(i.id) as purchase_orders_count 
	FROM buyers e 
	LEFT JOIN purchase_orders i ON i.buyer_id = e.id 
	WHERE e.id = ?
	GROUP BY e.id;`

	row := r.db.QueryRow(query, id)
	e := CountByBuyer{}
	err := row.Scan(&e.ID, &e.CardNumberID, &e.FirstName, &e.LastName, &e.Count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CountByBuyer{}, ErrNotFound
		}
		return CountByBuyer{}, ErrInternalServerError
	}

	return e, nil

}
