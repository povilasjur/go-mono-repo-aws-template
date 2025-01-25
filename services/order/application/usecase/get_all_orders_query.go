package usecase

import (
	"common"
	"context"
	"order/domain"
)

type GetAllOrdersQuery struct {
	Filter *domain.OrderFilter
	Page   *common.PageFilter
}

type GetAllOrdersQueryHandler struct {
	orderRepository domain.OrderRepository
}

func NewGetAllOrdersQueryHandler(orderRepository domain.OrderRepository) *GetAllOrdersQueryHandler {
	return &GetAllOrdersQueryHandler{
		orderRepository: orderRepository,
	}
}

func (h *GetAllOrdersQueryHandler) Execute(ctx context.Context, q GetAllOrdersQuery) (*common.Paginated[domain.Order], error) {
	return h.orderRepository.GetAll(ctx, q.Filter, q.Page)
}
