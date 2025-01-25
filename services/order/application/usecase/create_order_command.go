package usecase

import (
	"context"
	"order/domain"
)

type CreateOrderCommand struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateOrderCommandHandler struct {
	orderRepository domain.OrderRepository
}

func NewCreateOrderCommandHandler(orderRepository domain.OrderRepository) *CreateOrderCommandHandler {
	return &CreateOrderCommandHandler{
		orderRepository: orderRepository,
	}
}

func (h *CreateOrderCommandHandler) Execute(ctx context.Context, cmd CreateOrderCommand) (*domain.Order, error) {
	order, err := domain.CreateOrder(
		ctx,
		cmd.Id,
		cmd.Name,
	)
	if err != nil {
		return nil, err
	}
	return order, h.orderRepository.Save(ctx, order)
}
