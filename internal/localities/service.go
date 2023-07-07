package localities

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/optional"
)

type CreateDTO struct {
	Name     string
	Province string
	Country  string
}

type CountByLocality struct {
	ID    int
	Name  string
	Count int
}

type Service interface {
	Create(c context.Context, loc CreateDTO) (domain.Locality, error)
	// godoc CountSellers
	//  Returns slice of Seller count by Locality. If id is specified,
	//  the slice will only contain the count for that locality; otherwise,
	//  the slice will contain all localities.
	CountSellers(c context.Context, id optional.Opt[int]) ([]CountByLocality, error)
	// godoc CountCarriers
	//  Returns slice of Carrier count by Locality. If id is specified,
	//  the slice will only contain the count for that locality; otherwise,
	//  the slice will contain all localities.
	CountCarriers(c context.Context, id optional.Opt[int]) ([]CountByLocality, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (svc *service) Create(c context.Context, loc CreateDTO) (domain.Locality, error) {
	locDomain := MapCreateToDomain(loc)

	id, err := svc.repo.Save(c, locDomain)
	if err != nil {
		var errInvalidLoc *ErrInvalidLocality
		if errors.As(err, &errInvalidLoc) {
			return domain.Locality{}, err
		} else {
			return domain.Locality{}, NewErrGeneric("error saving locality")
		}
	}
	locDomain.ID = id

	return locDomain, nil
}

func (svc *service) CountSellers(c context.Context, id optional.Opt[int]) ([]CountByLocality, error) {
	locs, err := svc.repo.GetAll(c)
	if err != nil {
		return nil, NewErrGeneric("error fetching localities")
	}

	index := getLocalityIndex(locs)
	ids := getReportIDs(id, index)

	stats, err := svc.repo.CountSellersByLocalities(c, ids)
	if err != nil {
		return nil, NewErrGeneric("error counting sellers")
	}
	if id.HasVal && len(stats) == 0 {
		return nil, NewErrNotFound(id.Val)
	}

	report := make([]CountByLocality, 0)
	for _, stat := range stats {
		loc := index[stat.LocalityID]
		locCount := newSellersByLocality(loc, stat.Count)
		report = append(report, locCount)
	}

	return report, nil
}

func (svc *service) CountCarriers(c context.Context, id optional.Opt[int]) ([]CountByLocality, error) {
	locs, err := svc.repo.GetAll(c)
	if err != nil {
		return nil, NewErrGeneric("error fetching localities")
	}

	index := getLocalityIndex(locs)
	ids := getReportIDs(id, index)

	stats, err := svc.repo.CountCarriersByLocalities(c, ids)
	if err != nil {
		return nil, NewErrGeneric("error counting carriers")
	}
	if id.HasVal && len(stats) == 0 {
		return nil, NewErrNotFound(id.Val)
	}

	report := make([]CountByLocality, 0)
	for _, stat := range stats {
		loc := index[stat.LocalityID]
		locCount := newSellersByLocality(loc, stat.Count)
		report = append(report, locCount)
	}

	return report, nil
}

func newSellersByLocality(loc domain.Locality, count int) CountByLocality {
	return CountByLocality{
		ID:    loc.ID,
		Name:  loc.Name,
		Count: count,
	}
}

func getReportIDs(id optional.Opt[int], idx map[int]domain.Locality) []int {
	if ID, hasVal := id.Value(); hasVal {
		return []int{ID}
	}

	ids := make([]int, 0)
	for _, loc := range idx {
		ids = append(ids, loc.ID)
	}
	return ids
}

func getLocalityIndex(locs []domain.Locality) map[int]domain.Locality {
	index := make(map[int]domain.Locality)
	for _, loc := range locs {
		index[loc.ID] = loc
	}
	return index
}

func MapCreateToDomain(dto CreateDTO) domain.Locality {
	return domain.Locality{
		Name:     dto.Name,
		Province: dto.Province,
		Country:  dto.Country,
	}
}
