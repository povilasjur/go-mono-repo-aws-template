package usecase

import (
	apperrors "common/errors"
	"context"
	"order/domain"
)

type GetOrderQuery struct {
	Id string `json:"id"`
}

type GetOrderQueryHandler struct {
	orderRepository domain.OrderRepository
}

func NewGetOrderQueryHandler(orderRepository domain.OrderRepository) *GetOrderQueryHandler {
	return &GetOrderQueryHandler{
		orderRepository: orderRepository,
	}
}

func (h *GetOrderQueryHandler) Execute(ctx context.Context, query GetOrderQuery) (*domain.Order, error) {
	if len(query.Id) < 1 {
		return nil, apperrors.InvalidRequestParameter("id can not be empty", "id")
	}
	return h.orderRepository.GetById(ctx, query.Id)
}
