package application

import "order/application/usecase"

type OrderApplication struct {
	GetOrderQueryHandler      *usecase.GetOrderQueryHandler
	GetAllOrdersQueryHandler  *usecase.GetAllOrdersQueryHandler
	CreateOrderCommandHandler *usecase.CreateOrderCommandHandler
}

func NewOrderApplication(
	getOrderQueryHandler *usecase.GetOrderQueryHandler,
	getAllOrdersQueryHandler *usecase.GetAllOrdersQueryHandler,
	createOrderCommandHandler *usecase.CreateOrderCommandHandler,
) *OrderApplication {
	return &OrderApplication{
		GetOrderQueryHandler:      getOrderQueryHandler,
		GetAllOrdersQueryHandler:  getAllOrdersQueryHandler,
		CreateOrderCommandHandler: createOrderCommandHandler,
	}
}
