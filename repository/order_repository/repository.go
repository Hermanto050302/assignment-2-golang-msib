package order_repository

import (
	"assignment-2-golang-msib/entity"
	"net/http"
)

type OrderRepository interface {
	CreateOrder(orderPayload entity.Order, itemsPayload []entity.Item) (*entity.Order, error)
	UpdateOrder(orderPayload entity.Order, itemsPayload []entity.Item) (*OrderItem, error)
	GetAllOrders() ([]OrderItem, error)
	DeleteOrder(orderId int) (*http.Response, error)
}
