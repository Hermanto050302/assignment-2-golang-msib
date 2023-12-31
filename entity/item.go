package entity

import (
	"assignment-2-golang-msib/dto"
	"time"
)

type Item struct {
	ItemId      int
	ItemCode    string
	Quantity    int
	Description string
	OrderId     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (i *Item) ItemToItemResponse() dto.ItemResponse {
	return dto.ItemResponse{
		Id:          i.ItemId,
		ItemCode:    i.ItemCode,
		Quantity:    i.Quantity,
		Description: i.Description,
		OrderId:     i.OrderId,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}
