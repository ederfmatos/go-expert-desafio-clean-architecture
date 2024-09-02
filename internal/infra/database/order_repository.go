package database

import (
	"database/sql"
	"go-expert-desafio-clean-architecture/internal/entity"
	"go-expert-desafio-clean-architecture/internal/repository"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) List() ([]entity.Order, error) {
	orders := make([]entity.Order, 0)
	row, err := r.Db.Query("select id, price, tax, final_price from orders")
	if err != nil {
		return orders, err
	}
	for row.Next() {
		order := entity.Order{}
		err = row.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
