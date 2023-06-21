package stubs

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

type RepositoryStub struct{}

func (r *RepositoryStub) GetAll(ctx context.Context) ([]domain.Seller, error) {
	return []domain.Seller{}, nil
}

func (r *RepositoryStub) Get(ctx context.Context, id int) (domain.Seller, error) {
	return domain.Seller{}, nil
}

func (r *RepositoryStub) Exists(ctx context.Context, cid int) bool {
	return false
}

func (r *RepositoryStub) Save(ctx context.Context, s domain.Seller) (int, error) {
	return 1, nil
}

func (r *RepositoryStub) Update(ctx context.Context, s domain.Seller) error {
	return nil
}

func (r *RepositoryStub) Delete(ctx context.Context, id int) error {
	return nil
}
