package api

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/core/review"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/core/review/models"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/repository"
	pb "gitlab.ozon.dev/Woofka/movie-review-system-db/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(review review.Interface) pb.RepositoryServer {
	return &implementation{
		review: review,
	}
}

type implementation struct {
	pb.UnimplementedRepositoryServer
	review review.Interface
}

func (i *implementation) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	log.Print("API call of CreateReview")
	err := i.review.Create(ctx, &models.Review{
		Reviewer:   req.Review.GetReviewer(),
		MovieTitle: req.Review.GetMovieTitle(),
		Text:       req.Review.GetText(),
		Rating:     uint8(req.Review.GetRating()),
	})
	if err != nil {
		log.Print(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateReviewResponse{}, nil
}

func (i *implementation) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewResponse, error) {
	log.Print("API call of GetReview")
	r, err := i.review.Get(ctx, uint(req.GetId()))
	if err != nil {
		log.Print(err.Error())
		if errors.Is(err, repository.ErrReviewNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetReviewResponse{
		Review: &pb.Review{
			Id:         uint64(r.Id),
			Reviewer:   r.Reviewer,
			MovieTitle: r.MovieTitle,
			Text:       r.Text,
			Rating:     uint32(r.Rating),
		},
	}, nil
}

func (i *implementation) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.UpdateReviewResponse, error) {
	log.Print("API call of UpdateReview")
	err := i.review.Update(ctx, &models.Review{
		Id:         uint(req.Review.GetId()),
		Reviewer:   req.Review.GetReviewer(),
		MovieTitle: req.Review.GetMovieTitle(),
		Text:       req.Review.GetText(),
		Rating:     uint8(req.Review.GetRating()),
	})
	if err != nil {
		log.Print(err.Error())
		if errors.Is(err, repository.ErrReviewNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateReviewResponse{}, nil
}

func (i *implementation) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewResponse, error) {
	log.Print("API call of DeleteReview")
	err := i.review.Delete(ctx, uint(req.GetId()))
	if err != nil {
		log.Print(err.Error())
		if errors.Is(err, repository.ErrReviewNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteReviewResponse{}, nil
}

func (i *implementation) ListReview(ctx context.Context, req *pb.ListReviewRequest) (*pb.ListReviewResponse, error) {
	log.Print("API call of ListReview")
	reviews, err := i.review.List(ctx, uint(req.GetLimit()), uint(req.GetOffset()), req.GetOrderDesc())
	if err != nil {
		log.Print(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.Review, 0, len(reviews))
	for _, r := range reviews {
		result = append(result, &pb.Review{
			Id:         uint64(r.Id),
			Reviewer:   r.Reviewer,
			MovieTitle: r.MovieTitle,
			Text:       r.Text,
			Rating:     uint32(r.Rating),
		})
	}

	return &pb.ListReviewResponse{
		Reviews: result,
	}, nil
}
