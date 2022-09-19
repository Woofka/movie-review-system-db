package main

import (
	"net"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/api"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/core/review"
	pb "gitlab.ozon.dev/Woofka/movie-review-system-db/pkg/api"
	"google.golang.org/grpc"
)

func runGRPCServer(review review.Interface) error {
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		return errors.Wrap(err, "GRPC server listener error")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRepositoryServer(grpcServer, api.New(review))

	err = grpcServer.Serve(listener)
	if err != nil {
		return errors.Wrap(err, "GRPC server error")
	}

	return nil
}
