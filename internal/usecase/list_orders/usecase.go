package list_orders

import (
	"go-expert-desafio-clean-architecture/internal/repository"
)

type (
	Output struct {
		ID         string  `json:"id"`
		Price      float64 `json:"price"`
		Tax        float64 `json:"tax"`
		FinalPrice float64 `json:"final_price"`
	}

	UseCase struct {
		OrderRepository repository.OrderRepository
	}
)

func New(OrderRepository repository.OrderRepository) *UseCase {
	return &UseCase{OrderRepository: OrderRepository}
}

func (c *UseCase) Execute() (*[]Output, error) {
	orders, err := c.OrderRepository.List()
	if err != nil {
		return nil, err
	}
	output := make([]Output, len(orders))
	for i, order := range orders {
		output[i] = Output{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		}
	}
	return &output, nil
}
