package web

import (
	"encoding/json"
	"go-expert-desafio-clean-architecture/internal/repository"
	"go-expert-desafio-clean-architecture/internal/usecase/create_order"
	"go-expert-desafio-clean-architecture/internal/usecase/list_orders"
	"go-expert-desafio-clean-architecture/pkg/events"
	"net/http"
)

type OrderHandler struct {
	EventDispatcher   events.EventDispatcher
	OrderRepository   repository.OrderRepository
	OrderCreatedEvent events.Event
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcher,
	OrderRepository repository.OrderRepository,
	OrderCreatedEvent events.Event,
) *OrderHandler {
	return &OrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input create_order.Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createOrder := create_order.New(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *OrderHandler) List(w http.ResponseWriter, _ *http.Request) {
	listOrders := list_orders.New(h.OrderRepository)
	output, err := listOrders.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
