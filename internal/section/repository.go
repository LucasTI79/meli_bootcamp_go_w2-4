package section

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Repository encapsulates the storage of a section.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Section, error)
	Get(ctx context.Context, id int) (domain.Section, error)
	Exists(ctx context.Context, sectionNumber int) bool
	Save(ctx context.Context, s domain.Section) (int, error)
	Update(ctx context.Context, s domain.Section) error
	Delete(ctx context.Context, id int) error
	GetAllReportProducts(ctx context.Context) ([]domain.GetOneData, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Section, error) {
	query := "SELECT * FROM sections;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var sections []domain.Section

	for rows.Next() {
		s := domain.Section{}
		_ = rows.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
		sections = append(sections, s)
	}

	return sections, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Section, error) {
	query := "SELECT * FROM sections WHERE id=?;"
	row := r.db.QueryRow(query, id)
	s := domain.Section{}
	err := row.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Section{}, ErrNotFound
		}
		return domain.Section{}, err
	}

	return s, nil
}

func (r *repository) Exists(ctx context.Context, sectionNumber int) bool {
	query := "SELECT section_number FROM sections WHERE section_number=?;"
	row := r.db.QueryRow(query, sectionNumber)
	err := row.Scan(&sectionNumber)
	return err == nil
}

func (r *repository) Save(ctx context.Context, s domain.Section) (int, error) {
	query := `INSERT INTO sections
		(section_number, current_temperature, minimum_temperature,
		current_capacity, minimum_capacity, maximum_capacity,
		warehouse_id, product_type_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, s domain.Section) error {
	query := `UPDATE sections SET section_number=?, current_temperature=?,
		minimum_temperature=?, current_capacity=?, minimum_capacity=?,
		maximum_capacity=?, warehouse_id=?, product_type_id=?
		WHERE id=?;`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID, &s.ID)
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
	query := "DELETE FROM sections WHERE id=?;"
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

func (r *repository) GetAllReportProducts(ctx context.Context) ([]domain.GetOneData, error) {
	var sections []domain.GetOneData
	query := "SELECT s.id, s.section_number, COUNT(p.id) AS products_count FROM sections s INNER JOIN product_batches pb ON s.ID = pb.section_id INNER JOIN products p ON pb.product_id = p.id GROUP by s.id, s.section_number;"
	rows, err := r.db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return sections, err
	}
	for rows.Next() {
		s := domain.GetOneData{}
		err = rows.Scan(&s.SectionId, &s.SectionNumber, &s.ProductCount)
		if err != nil {
			fmt.Println(err.Error())
			return sections, err
		}
		sections = append(sections, s)
	}
	return sections, nil
}
