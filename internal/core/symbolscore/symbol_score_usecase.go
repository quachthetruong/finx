package symbolscore

import (
	"context"
	"fmt"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	stockExchangeRepo "financing-offer/internal/core/stockexchange/repository"
	"financing-offer/internal/core/symbolscore/repository"
)

type UseCase interface {
	GetAll(ctx context.Context, filter entity.SymbolScoreFilter) ([]entity.SymbolScore, error)
	Update(ctx context.Context, symbolScore entity.SymbolScore) (entity.SymbolScore, error)
	Create(ctx context.Context, symbolScore entity.SymbolScore) (entity.SymbolScore, error)
	GetById(ctx context.Context, id int64) (entity.SymbolScore, error)
}

func NewUseCase(
	symbolScoreRepo repository.SymbolScoreRepository,
	stockExchangeRepository stockExchangeRepo.StockExchangeRepository,
) UseCase {
	return &symbolScoreUseCase{repository: symbolScoreRepo, stockExchangeRepository: stockExchangeRepository}
}

type symbolScoreUseCase struct {
	repository              repository.SymbolScoreRepository
	stockExchangeRepository stockExchangeRepo.StockExchangeRepository
}

func (s *symbolScoreUseCase) GetById(ctx context.Context, id int64) (entity.SymbolScore, error) {
	symbolScore, err := s.repository.GetById(ctx, id)
	if err != nil {
		return symbolScore, fmt.Errorf("symbolScoreUseCase GetById %w", err)
	}
	return symbolScore, nil
}

func (s *symbolScoreUseCase) GetAll(ctx context.Context, filter entity.SymbolScoreFilter) ([]entity.SymbolScore, error) {
	symbolScores, err := s.repository.GetAll(ctx, filter)
	if err != nil {
		return symbolScores, fmt.Errorf("symbolScoreUseCase GetAll %w", err)
	}
	return symbolScores, nil
}

func (s *symbolScoreUseCase) Update(ctx context.Context, symbolScore entity.SymbolScore) (entity.SymbolScore, error) {
	updated, err := s.repository.Update(ctx, symbolScore)
	if err != nil {
		return updated, fmt.Errorf("symbolScoreUseCase Update %w", err)
	}
	return updated, nil
}

func (s *symbolScoreUseCase) Create(ctx context.Context, symbolScore entity.SymbolScore) (entity.SymbolScore, error) {
	stockExchange, err := s.stockExchangeRepository.GetBySymbolId(ctx, symbolScore.SymbolId)
	if err != nil {
		return symbolScore, fmt.Errorf("symbolScoreUseCase Create %w", err)
	}
	if symbolScore.Score > stockExchange.MaxScore || symbolScore.Score < stockExchange.MinScore {
		return symbolScore, apperrors.ErrInvalidSymbolScore
	}
	created, err := s.repository.Create(ctx, symbolScore)
	if err != nil {
		return created, fmt.Errorf("symbolScoreUseCase Create %w", err)
	}
	return created, nil
}
