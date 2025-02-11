package services

import (
	"context"
	"fmt"
	"log/slog"
)

type hello interface {
	Hello(ctx context.Context) (string, error)
}

type HelloService struct {
	log   *slog.Logger
	hello hello
}

func NewHelloService(
	log *slog.Logger,
	hello hello,
) *HelloService {
	return &HelloService{
		log:   log,
		hello: hello,
	}
}

func (s *HelloService) Hello(
	ctx context.Context,
) (string, error) {
	const op = "helloService.Hello"
	log := s.log.With(slog.String("op", op))

	log.Info("retrieving hello")

	res, err := s.hello.Hello(ctx)
	if err != nil {
		log.Error("failed to retrieve cache entries")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("hello retrieved")

	return res, nil
}
