package stockexchange

import (
	"context"
	"fmt"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/stockexchange/repository"
)

type UseCase interface {
	GetAll(ctx context.Context) ([]entity.StockExchange, error)
	Create(ctx context.Context, stockExchange entity.StockExchange) (entity.StockExchange, error)
	Update(ctx context.Context, stockExchange entity.StockExchange) (entity.StockExchange, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (entity.StockExchange, error)
}

type stockExchangeUseCase struct {
	repository repository.StockExchangeRepository
}

func (s *stockExchangeUseCase) GetById(ctx context.Context, id int64) (entity.StockExchange, error) {
	res, err := s.repository.GetById(ctx, id)
	if err != nil {
		return res, fmt.Errorf("stockExchangeUseCase GetById %w", err)
	}
	return res, nil
}

func NewUseCase(stockExchangeRepo repository.StockExchangeRepository) UseCase {
	return &stockExchangeUseCase{repository: stockExchangeRepo}
}

func (s *stockExchangeUseCase) GetAll(ctx context.Context) ([]entity.StockExchange, error) {
	res, err := s.repository.GetAll(ctx)
	if err != nil {
		return res, fmt.Errorf("stockExchangeUseCase GetAll %w", err)
	}
	return res, nil
}

func (s *stockExchangeUseCase) Create(ctx context.Context, stockExchange entity.StockExchange) (entity.StockExchange, error) {
	res, err := s.repository.Create(ctx, stockExchange)
	if err != nil {
		return res, fmt.Errorf("stockExchangeUseCase Create %w", err)
	}
	return res, err
}

func (s *stockExchangeUseCase) Update(ctx context.Context, stockExchange entity.StockExchange) (entity.StockExchange, error) {
	res, err := s.repository.Update(ctx, stockExchange)
	if err != nil {
		return res, fmt.Errorf("stockExchangeUseCase Update %w", err)
	}
	return res, nil
}

func (s *stockExchangeUseCase) Delete(ctx context.Context, id int64) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("stockExchangeUseCase Delete %w", err)
	}
	return nil
}
