package symbol

import (
	"errors"
	"fmt"
	"sort"

	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/atomicity"
	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/symbol/repository"
	symbolScoreRepo "financing-offer/internal/core/symbolscore/repository"
	"financing-offer/internal/funcs"
	"financing-offer/pkg/optional"
)

type UseCase interface {
	GetAll(ctx context.Context, filter entity.SymbolFilter) ([]entity.Symbol, core.PagingMetaData, error)
	UpdateStatus(ctx context.Context, symbol entity.Symbol) (entity.Symbol, error)
	Create(ctx context.Context, symbol entity.Symbol) (entity.Symbol, error)
	GetById(ctx context.Context, id int64) (entity.Symbol, error)
	GetBySymbol(ctx context.Context, symbol string) (entity.Symbol, error)
	GetSymbolNotInActiveBlacklist(ctx context.Context, symbol string) (entity.Symbol, error)
}

type symbolUseCase struct {
	repository            repository.SymbolRepository
	symbolScoreRepository symbolScoreRepo.SymbolScoreRepository
	atomicExecutor        atomicity.AtomicExecutor
}

func NewUseCase(
	symbolRepo repository.SymbolRepository,
	symbolScoreRepository symbolScoreRepo.SymbolScoreRepository,
	atomicExecutor atomicity.AtomicExecutor,
) UseCase {
	return &symbolUseCase{
		repository:            symbolRepo,
		symbolScoreRepository: symbolScoreRepository,
		atomicExecutor:        atomicExecutor,
	}
}

func (s *symbolUseCase) GetAll(ctx context.Context, filter entity.SymbolFilter) ([]entity.Symbol, core.PagingMetaData, error) {
	var (
		eg              errgroup.Group
		symbolsEntities []entity.Symbol
		pagingMetaData  = core.PagingMetaData{PageSize: filter.Size, PageNumber: filter.Number}
	)
	eg.Go(
		func() error {
			res, scopedErr := s.repository.GetAll(ctx, filter)
			symbolsEntities = res
			return scopedErr
		},
	)
	eg.Go(
		func() error {
			res, scopedErr := s.repository.Count(ctx, filter)
			pagingMetaData.Total = res
			pagingMetaData.TotalPages = filter.TotalPages(res)
			return scopedErr
		},
	)
	if err := eg.Wait(); err != nil {
		return symbolsEntities, pagingMetaData, fmt.Errorf("symbolUseCase GetAll %w", err)
	}
	symbols := funcs.Map(
		symbolsEntities, func(s entity.Symbol) string {
			return s.Symbol
		},
	)
	symbolScores, err := s.symbolScoreRepository.GetAll(
		ctx, entity.SymbolScoreFilter{
			Symbols: symbols,
		},
	)
	if err != nil {
		return symbolsEntities, pagingMetaData, fmt.Errorf("symbolUseCase GetAll %w", err)
	}
	scoreBySymbolId := make(map[int64][]entity.SymbolScore)
	for _, symbolScore := range symbolScores {
		currentScores := scoreBySymbolId[symbolScore.SymbolId]
		scoreBySymbolId[symbolScore.SymbolId] = append(currentScores, symbolScore)
	}
	res := funcs.Map(
		symbolsEntities, func(symbol entity.Symbol) entity.Symbol {
			scores := scoreBySymbolId[symbol.Id]
			sort.SliceStable(
				scores, func(i, j int) bool {
					return scores[i].AffectedFrom.Before(scores[j].AffectedFrom)
				},
			)
			if scores == nil {
				scores = []entity.SymbolScore{}
			}
			symbol.Scores = scores
			return symbol
		},
	)
	return res, pagingMetaData, nil
}

func (s *symbolUseCase) UpdateStatus(ctx context.Context, newSymbol entity.Symbol) (entity.Symbol, error) {
	errorWrapMsg := "symbolUseCase UpdateStatus %w"
	res, err := s.repository.Update(ctx, newSymbol)
	if err != nil {
		return entity.Symbol{}, fmt.Errorf(errorWrapMsg, err)
	}
	return res, nil
}

func (s *symbolUseCase) Create(ctx context.Context, symbol entity.Symbol) (entity.Symbol, error) {
	res, err := s.repository.Create(ctx, symbol)
	if err != nil {
		return res, fmt.Errorf("symbolUseCase Create %w", err)
	}
	return res, nil
}

func (s *symbolUseCase) GetById(ctx context.Context, id int64) (entity.Symbol, error) {
	symbolEntity, err := s.repository.GetById(ctx, id)
	if err != nil {
		return symbolEntity, fmt.Errorf("symbolUseCase GetById %w", err)
	}
	scores, err := s.symbolScoreRepository.GetAll(
		ctx, entity.SymbolScoreFilter{
			Symbols: []string{symbolEntity.Symbol},
			Status:  optional.Some(entity.SymbolScoreStatusActive),
		},
	)
	if err != nil {
		return symbolEntity, fmt.Errorf("symbolUseCase GetById %w", err)
	}
	sort.SliceStable(
		scores, func(i, j int) bool {
			return scores[i].AffectedFrom.Before(scores[j].AffectedFrom)
		},
	)
	symbolEntity.Scores = scores
	return symbolEntity, nil
}

func (s *symbolUseCase) GetBySymbol(ctx context.Context, symbol string) (entity.Symbol, error) {
	symbolEntity, err := s.repository.GetBySymbol(ctx, symbol)
	if err != nil {
		return symbolEntity, fmt.Errorf("symbolUseCase GetBySymbol %w", err)
	}
	scores, err := s.symbolScoreRepository.GetAll(
		ctx, entity.SymbolScoreFilter{
			Symbols: []string{symbolEntity.Symbol},
			Status:  optional.Some(entity.SymbolScoreStatusActive),
		},
	)
	if err != nil {
		return symbolEntity, fmt.Errorf("symbolUseCase GetBySymbol %w", err)
	}
	sort.SliceStable(
		scores, func(i, j int) bool {
			return scores[i].AffectedFrom.Before(scores[j].AffectedFrom)
		},
	)
	symbolEntity.Scores = scores
	return symbolEntity, nil
}

func (s *symbolUseCase) GetSymbolNotInActiveBlacklist(ctx context.Context, symbol string) (entity.Symbol, error) {
	symbolEntity, err := s.repository.GetBySymbol(ctx, symbol)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return entity.Symbol{}, apperrors.ErrSymbolCodeNotFound
		}
		return entity.Symbol{}, fmt.Errorf("symbolUseCase GetSymbolNotInActiveBlacklist %w", err)
	}

	_, err = s.repository.GetSymbolWithActiveBlacklist(ctx, symbol)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return symbolEntity, nil
		}
		return entity.Symbol{}, fmt.Errorf("symbolUseCase GetSymbolNotInActiveBlacklist %w", err)
	}
	return entity.Symbol{}, apperrors.ErrSymbolCodeInBlacklist
}
