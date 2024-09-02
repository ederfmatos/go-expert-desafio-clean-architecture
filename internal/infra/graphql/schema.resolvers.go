package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"go-expert-desafio-clean-architecture/internal/infra/graphql/model"
	"go-expert-desafio-clean-architecture/internal/usecase/create_order"
)

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input *model.OrderInput) (*model.Order, error) {
	dto := create_order.Input{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	output, err := r.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &model.Order{
		ID:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

// ListOrders is the resolver for the listOrders field.
func (r *queryResolver) ListOrders(ctx context.Context) ([]*model.Order, error) {
	output, err := r.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}
	orders := make([]*model.Order, len(*output))
	for i, order := range *output {
		orders[i] = &model.Order{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		}
	}
	return orders, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }