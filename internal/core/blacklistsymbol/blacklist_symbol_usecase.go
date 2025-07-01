package blacklistsymbol

import (
	"context"
	"fmt"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/blacklistsymbol/repository"
	"financing-offer/internal/core/entity"
)

type UseCase interface {
	GetAll(ctx context.Context, filter entity.BlacklistSymbolFilter) ([]entity.BlacklistSymbol, error)
	Update(ctx context.Context, symbol entity.BlacklistSymbol) (entity.BlacklistSymbol, error)
	Create(ctx context.Context, symbol entity.BlacklistSymbol) (entity.BlacklistSymbol, error)
	GetById(ctx context.Context, id int64) (entity.BlacklistSymbol, error)
}

func NewUseCase(blacklistSymbolRepo repository.BlackListSymbolRepository) UseCase {
	return &blacklistSymbolUseCase{repository: blacklistSymbolRepo}
}

type blacklistSymbolUseCase struct {
	repository repository.BlackListSymbolRepository
}

func (u *blacklistSymbolUseCase) GetAll(ctx context.Context, filter entity.BlacklistSymbolFilter) ([]entity.BlacklistSymbol, error) {
	return u.repository.GetAll(ctx, filter)
}

func (u *blacklistSymbolUseCase) Update(ctx context.Context, blacklistSymbol entity.BlacklistSymbol) (entity.BlacklistSymbol, error) {
	if blacklistSymbol.Status == entity.BlacklistSymbolStatusActive {
		overlap, err := u.repository.GetByAffectTime(ctx, blacklistSymbol.SymbolId, blacklistSymbol.AffectedFrom, blacklistSymbol.AffectedTo)
		if err != nil {
			return entity.BlacklistSymbol{}, fmt.Errorf("BlackListSymbolUseCase Update %w", err)
		}
		// If more then one overlap, it means that the new time range overlaps with more than one record
		if len(overlap) > 1 {
			return entity.BlacklistSymbol{}, apperrors.ErrBlacklistSymbolOverlap
		}
		// If only one overlap, it means that the new time range overlaps with only one record, check if it is the same record
		if len(overlap) == 1 && overlap[0].Id != blacklistSymbol.Id {
			return entity.BlacklistSymbol{}, apperrors.ErrBlacklistSymbolOverlap
		}
	}
	return u.repository.Update(ctx, blacklistSymbol)
}

func (u *blacklistSymbolUseCase) Create(ctx context.Context, blacklistSymbol entity.BlacklistSymbol) (entity.BlacklistSymbol, error) {
	if blacklistSymbol.Status == entity.BlacklistSymbolStatusActive {
		overlap, err := u.repository.GetByAffectTime(ctx, blacklistSymbol.SymbolId, blacklistSymbol.AffectedFrom, blacklistSymbol.AffectedTo)
		if err != nil {
			return entity.BlacklistSymbol{}, fmt.Errorf("BlackListSymbolUseCase Create %w", err)
		}
		if len(overlap) > 0 {
			return entity.BlacklistSymbol{}, apperrors.ErrBlacklistSymbolOverlap
		}
	}
	return u.repository.Create(ctx, blacklistSymbol)
}

func (u *blacklistSymbolUseCase) GetById(ctx context.Context, id int64) (entity.BlacklistSymbol, error) {
	return u.repository.GetById(ctx, id)
}
