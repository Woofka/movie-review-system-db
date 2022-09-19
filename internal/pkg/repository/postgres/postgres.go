package postgres

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/core/review/models"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/repository"
)

func New(pool *pgxpool.Pool) repository.Interface {
	return &Repository{pool}
}

type Repository struct {
	pool *pgxpool.Pool
}

func (r *Repository) List(ctx context.Context, limit, offset uint, orderDesc bool) ([]*models.Review, error) {
	query := "SELECT id, reviewer, movie_title, text, rating FROM public.reviews ORDER BY ID "
	if orderDesc {
		query += " DESC "
	}
	query += "LIMIT $1 OFFSET $2"
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "postgres.List")
	}
	defer rows.Close()

	var reviews []*models.Review
	err = pgxscan.ScanAll(&reviews, rows)
	if err != nil {
		return nil, errors.Wrap(err, "postgres.List")
	}

	return reviews, nil
}

func (r *Repository) addUser(ctx context.Context, username string) error {
	query := "INSERT INTO public.users (username) values ($1) on conflict do nothing"
	_, err := r.pool.Exec(ctx, query, username)
	if err != nil {
		return errors.Wrap(err, "postgres.addUser")
	}
	return nil
}

func (r *Repository) addMovie(ctx context.Context, movieTitle string) error {
	query := "INSERT INTO public.movies (title) values ($1) on conflict do nothing"
	_, err := r.pool.Exec(ctx, query, movieTitle)
	if err != nil {
		return errors.Wrap(err, "postgres.addMovie")
	}
	return nil
}

func (r *Repository) Create(ctx context.Context, review *models.Review) error {
	err := r.addUser(ctx, review.Reviewer)
	if err != nil {
		return errors.Wrap(err, "postgres.Create")
	}

	err = r.addMovie(ctx, review.MovieTitle)
	if err != nil {
		return errors.Wrap(err, "postgres.Create")
	}

	query := "INSERT INTO public.reviews (reviewer, movie_title, text, rating) " +
		"values ($1, $2, $3, $4) RETURNING id"

	row := r.pool.QueryRow(ctx, query, review.Reviewer, review.MovieTitle, review.Text, int(review.Rating))
	var reviewId int
	err = row.Scan(&reviewId)
	if err != nil {
		return errors.Wrap(err, "postgres.Create")
	}

	review.Id = uint(reviewId)
	return nil
}

func (r *Repository) Get(ctx context.Context, id uint) (*models.Review, error) {
	query := "SELECT id, reviewer, movie_title, text, rating FROM public.reviews WHERE id = $1"
	row := r.pool.QueryRow(ctx, query, id)
	review := models.Review{}
	err := row.Scan(&review.Id, &review.Reviewer, &review.MovieTitle, &review.Text, &review.Rating)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrapf(repository.ErrReviewNotExists, "review with id %d does not exists", id)
		}
		return nil, errors.Wrap(err, "postgres.Get")
	}

	return &review, nil
}

func (r *Repository) Update(ctx context.Context, review *models.Review) error {
	query := "UPDATE public.reviews SET text = $1, rating = $2 WHERE id = $3"
	result, err := r.pool.Exec(ctx, query, review.Text, review.Rating, review.Id)
	if err != nil {
		return errors.Wrap(err, "postgres.Update")
	}
	if result.RowsAffected() == 0 {
		return errors.Wrapf(repository.ErrReviewNotExists, "review with id %d does not exists", review.Id)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	query := "DELETE FROM public.reviews WHERE id = $1"
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "postgres.Delete")
	}
	if result.RowsAffected() == 0 {
		return errors.Wrapf(repository.ErrReviewNotExists, "review with id %d does not exists", id)
	}
	return nil
}
