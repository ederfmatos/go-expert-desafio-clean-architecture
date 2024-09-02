package service

import (
	"context"
	"go-expert-desafio-clean-architecture/internal/infra/grpc/pb"
	"go-expert-desafio-clean-architecture/internal/usecase/create_order"
	"go-expert-desafio-clean-architecture/internal/usecase/list_orders"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase *create_order.CreateOrderUseCase
	ListOrdersUseCase  *list_orders.UseCase
}

func NewOrderService(
	createOrderUseCase *create_order.CreateOrderUseCase,
	listOrdersUseCase *list_orders.UseCase,
) pb.OrderServiceServer {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(_ context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := create_order.Input{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(context.Context, *pb.Empty) (*pb.ListOrdersResponse, error) {
	output, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}
	response := pb.ListOrdersResponse{
		Orders: make([]*pb.ListOrdersItem, len(*output)),
	}
	for i, order := range *output {
		response.Orders[i] = &pb.ListOrdersItem{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.Price + order.Tax),
		}
	}
	return &response, nil
}
