package repository

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/core/review/models"
)

var ErrReviewNotExists = errors.New("review does not exist")

type Interface interface {
	Create(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, review *models.Review) error
	List(ctx context.Context, limit, offset uint, orderDesc bool) ([]*models.Review, error)
	Get(ctx context.Context, id uint) (*models.Review, error)
}
