package localities

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

type SellerCount struct {
	LocalityID  int
	SellerCount int
}

type Repository interface {
	Save(c context.Context, loc domain.Locality) (int, error)
	GetAll(c context.Context) ([]domain.Locality, error)
	CountSellersByLocalities(c context.Context, ids []int) ([]SellerCount, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(ctx context.Context, loc domain.Locality) (int, error) {
	// This query inserts the domain's locality across 3 tables.
	// It tries to insert the province and country names if they
	// don't exist, ignoring possible unique-constraint violations.
	// The last INSERT should not be ignored, since its failure means
	// that the whole locality already exists.
	countryQuery := `INSERT IGNORE INTO countries (country_name) VALUES (?);`
	provinceQuery := `
		INSERT IGNORE INTO provinces (province_name, country_id)
			SELECT ?, c.id FROM countries c
		WHERE c.country_name = ?;`
	localityQuery := `
		INSERT INTO localities (locality_name, province_id)
			SELECT ?, p.id FROM provinces p
		WHERE p.province_name = ?;`

	r.db.Exec(countryQuery, loc.Country)
	r.db.Exec(provinceQuery, loc.Province, loc.Country)
	result, err := r.db.Exec(localityQuery, loc.Name, loc.Province)

	if err != nil {
		if isDuplicateError(err) {
			return 0, NewErrInvalidLocality(loc)
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) GetAll(c context.Context) ([]domain.Locality, error) {
	query := `SELECT l.id, l.locality_name, p.province_name, c.country_name
		FROM countries c JOIN provinces p ON c.id = p.country_id
		JOIN localities l ON p.id = l.province_id;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	locs := make([]domain.Locality, 0)
	for rows.Next() {
		var loc domain.Locality
		err := rows.Scan(&loc.ID, &loc.Name, &loc.Province, &loc.Country)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		locs = append(locs, loc)
	}

	return locs, nil
}

func (r *repository) CountSellersByLocalities(c context.Context, ids []int) ([]SellerCount, error) {
	if len(ids) == 0 {
		return make([]SellerCount, 0), nil
	}

	query := `SELECT l.id, COUNT(s.id)
		FROM localities l
		LEFT JOIN sellers s ON s.locality_id = l.id
		WHERE l.id IN (?` + strings.Repeat(",?", len(ids)-1) + `)
		GROUP BY l.id, l.locality_name;`

	queryArgs := convertToAny(ids)
	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		return nil, err
	}

	counts := make([]SellerCount, 0)
	for rows.Next() {
		var count SellerCount
		err := rows.Scan(&count.LocalityID, &count.SellerCount)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		counts = append(counts, count)
	}

	return counts, nil
}

func isDuplicateError(err error) bool {
	return strings.HasPrefix(err.Error(), "Error 1062")
}

func convertToAny[T any](x []T) []any {
	ret := make([]any, len(x))

	for i, xi := range x {
		ret[i] = xi
	}

	return ret
}
