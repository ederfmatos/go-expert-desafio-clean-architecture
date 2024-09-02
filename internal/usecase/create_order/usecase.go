package create_order

import (
	"go-expert-desafio-clean-architecture/internal/entity"
	"go-expert-desafio-clean-architecture/internal/repository"
	"go-expert-desafio-clean-architecture/pkg/events"
)

type (
	Input struct {
		ID    string  `json:"id"`
		Price float64 `json:"price"`
		Tax   float64 `json:"tax"`
	}

	Output struct {
		ID         string  `json:"id"`
		Price      float64 `json:"price"`
		Tax        float64 `json:"tax"`
		FinalPrice float64 `json:"final_price"`
	}

	CreateOrderUseCase struct {
		OrderRepository repository.OrderRepository
		OrderCreated    events.Event
		EventDispatcher events.EventDispatcher
	}
)

func New(
	OrderRepository repository.OrderRepository,
	OrderCreated events.Event,
	EventDispatcher events.EventDispatcher,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input Input) (*Output, error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	err := order.CalculateFinalPrice()
	if err != nil {
		return nil, err
	}
	if err = c.OrderRepository.Save(&order); err != nil {
		return nil, err
	}
	output := Output{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}
	c.OrderCreated.SetPayload(output)
	err = c.EventDispatcher.Dispatch(c.OrderCreated)
	if err != nil {
		return nil, err
	}
	return &output, nil
}
