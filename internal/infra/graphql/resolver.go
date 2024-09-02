package graphql

import (
	"go-expert-desafio-clean-architecture/internal/usecase/create_order"
	"go-expert-desafio-clean-architecture/internal/usecase/list_orders"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateOrderUseCase *create_order.CreateOrderUseCase
	ListOrdersUseCase  *list_orders.UseCase
}
