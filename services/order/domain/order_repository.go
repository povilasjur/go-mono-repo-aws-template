package domain

import (
	"common"
	"context"
)

type OrderRepository interface {
	GetById(ctx context.Context, id string) (*Order, error)
	GetAll(ctx context.Context, merchantFilter *OrderFilter, pageFilter *common.PageFilter) (*common.Paginated[Order], error)
	Save(ctx context.Context, order *Order) error
}
