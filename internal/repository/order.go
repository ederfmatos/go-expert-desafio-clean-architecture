package repository

import "go-expert-desafio-clean-architecture/internal/entity"

type OrderRepository interface {
	Save(order *entity.Order) error
	List() ([]entity.Order, error)
}
