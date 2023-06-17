package order_repository

import "assignment-2-golang-msib/entity"

type OrderItem struct {
	Order entity.Order
	Items []entity.Item
}