package section

import (
	"context"
	"errors"
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/gin-gonic/gin"
)

// Errors
var (
	ErrNotFound = errors.New("section not found")
)

type Service interface {
	Save(ctx context.Context, section domain.CreateSection) (int, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Save(ctx context.Context, section domain.CreateSection) (int, error) {
	existsSectionNumber := s.repository.Exists(ctx, section.SectionNumber)
	if existsSectionNumber {
		 web.Error(&gin.Context{}, http.StatusConflict, "section number alredy exists")
		 return 0, nil
	}
	return s.repository.Save(ctx, section)
}
