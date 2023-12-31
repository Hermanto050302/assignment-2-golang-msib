package item_pg

import (
	"assignment-2-golang-msib/entity"
	"assignment-2-golang-msib/repository/item_repository"
	"database/sql"
	"errors"
	"fmt"
)

type itemPG struct {
	db *sql.DB
}

func NewItemPG(db *sql.DB) item_repository.ItemRepository {
	return &itemPG{
		db: db,
	}
}

func (i *itemPG) generatePlaceHolders(dataAmount int) string {
	start := "("

	for i := 1; i <= dataAmount; i++ {
		if i < dataAmount {
			start += fmt.Sprintf("$%d,", i)
		}

		if i == dataAmount {
			start += fmt.Sprintf("$%d)", i)
		}
	}

	return start
}

func (i *itemPG) findItemByItemCodesQuery(dataAmount int) string {
	query := `
		SELECT item_id, item_code, quantity, description, order_id, created_at, updated_at from "items"
		WHERE item_code IN 
	`
	placeHolders := i.generatePlaceHolders(dataAmount)

	result := query + placeHolders

	return result
}

func (i *itemPG) FindItemsByItemCodes(itemCodes []string) ([]*entity.Item, error) {
	query := i.findItemByItemCodesQuery(len(itemCodes))
	//SELECT * from "items" WHERE item_code IN ($1, $2)

	args := []any{}

	for _, value := range itemCodes {
		args = append(args, value)
	}

	rows, err := i.db.Query(query, args...)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer rows.Close()

	items := []*entity.Item{}
	for rows.Next() {
		item := entity.Item{}

		err = rows.Scan(&item.ItemId, &item.ItemCode, &item.Quantity, &item.Description, &item.OrderId, &item.CreatedAt, &item.UpdatedAt)

		if err != nil {
			return nil, errors.New(err.Error())
		}

		items = append(items, &item)
	}

	return items, nil
}
