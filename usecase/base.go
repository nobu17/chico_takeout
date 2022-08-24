package usecase

import (
	"context"

	"chico/takeout/common"
)

type BaseUseCase struct {
	ctx context.Context
}

func NewBaseUseCase() *BaseUseCase {
	return &BaseUseCase{}
}

func (b *BaseUseCase) InitContext(ctx context.Context) {
	b.ctx = ctx
}

func (b *BaseUseCase) IsAdmin() bool {
	return common.GetIsAdmin(b.ctx)
}
