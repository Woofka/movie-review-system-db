package review

import (
	"context"

	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/core/review/models"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/repository"
)

type Interface interface {
	Create(ctx context.Context, review *models.Review) error
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset uint, orderDesc bool) ([]*models.Review, error)
	Get(ctx context.Context, id uint) (*models.Review, error)
}

func New(repository repository.Interface) Interface {
	return &review{
		repository: repository,
	}
}

type review struct {
	repository repository.Interface
}

func (r *review) Create(ctx context.Context, review *models.Review) error {
	return r.repository.Create(ctx, review)
}

func (r *review) Update(ctx context.Context, review *models.Review) error {
	return r.repository.Update(ctx, review)
}

func (r *review) Delete(ctx context.Context, id uint) error {
	return r.repository.Delete(ctx, id)
}

func (r *review) List(ctx context.Context, limit, offset uint, orderDesc bool) ([]*models.Review, error) {
	return r.repository.List(ctx, limit, offset, orderDesc)
}

func (r *review) Get(ctx context.Context, id uint) (*models.Review, error) {
	return r.repository.Get(ctx, id)
}
