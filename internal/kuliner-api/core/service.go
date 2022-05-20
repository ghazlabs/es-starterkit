package core

import (
	"context"
	"fmt"

	"gopkg.in/validator.v2"
)

type Service interface {
	IndexFood(ctx context.Context, input FoodInput) (*Food, error)
	DeleteFood(ctx context.Context, id string) error
	SearchFoods(ctx context.Context, query string) ([]Food, error)
}

type Config struct {
	Storage     Storage `validate:"nonnil"`
	SearchLimit int
}

func NewService(cfg Config) (Service, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, err
	}
	searchLimit := cfg.SearchLimit
	if searchLimit <= 0 {
		searchLimit = defaultSearchLimit
	}
	s := &service{storage: cfg.Storage, searchLimit: searchLimit}
	return s, nil
}

type service struct {
	storage     Storage
	searchLimit int
}

func (s *service) IndexFood(ctx context.Context, input FoodInput) (*Food, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}
	err = s.storage.IndexFood(ctx, input.ToFood())
	if err != nil {
		return nil, fmt.Errorf("unable to index food into storage due: %w", err)
	}
	return nil, nil
}

func (s *service) DeleteFood(ctx context.Context, id string) error {
	if len(id) == 0 {
		return ErrMissingID
	}
	err := s.storage.DeleteFood(ctx, id)
	if err != nil {
		return fmt.Errorf("unable to delete food from storage due: %w", err)
	}
	return nil
}

func (s *service) SearchFoods(ctx context.Context, query string) ([]Food, error) {
	foods, err := s.storage.SearchFood(ctx, query, s.searchLimit)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("unable to search foods from storage due: %w", err)
	}
	return foods, nil
}
