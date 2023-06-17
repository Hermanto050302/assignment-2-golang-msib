package order_pg

import (
	"assignment-2-golang-msib/entity"
	"assignment-2-golang-msib/repository/order_repository"
	"database/sql"
	"errors"
	"net/http"
)

const (
	updateOrderQuery = `
		UPDATE "orders"
		SET ordered_at = $2,
		customer_name = $3
		WHERE order_id = $1
		RETURNING order_id, customer_name, created_at, updated_at
	`

	createOrderQuery = `
		INSERT INTO "orders"
			(
				ordered_at,
				customer_name
			)
		VALUES($1, $2)
		RETURNING order_id, customer_name, ordered_at, created_at,updated_at
	`
	createItemQuery = `
		INSERT INTO "items"
			(
				item_code,
				quantity,
				description,
				order_id
			)
		VALUES($1, $2, $3, $4)
		RETURNING item_id
	`

	updateItemQuery = `
		UPDATE "items"
		SET description = $2,
		quantity = $3
		WHERE item_code = $1
		RETURNING item_id, item_code, quantity, description, updated_at, order_id, created_at
	`

	getAllOrdersQuery = `
		SELECT order_id, customer_name, ordered_at, created_at, updated_at
		FROM "orders"
	`

	getItemsByOrderIdQuery = `
		SELECT item_id, item_code, quantity, description, created_at, updated_at, order_id
		FROM "items"
		WHERE order_id = $1
	`
	deleteOrderQuery = `
		DELETE FROM "orders"
		WHERE order_id = $1
	`
)

type orderPg struct {
	db *sql.DB
}

func NewOrderPG(db *sql.DB) order_repository.OrderRepository {
	return &orderPg{db: db}
}

//getItemsByOrderId
func (o *orderPg) getItemsByOrderId(orderId int) ([]entity.Item, error) {
	internalServerError := errors.New("something went wrong")

	rows, err := o.db.Query(getItemsByOrderIdQuery, orderId)

	if err != nil {
		return nil, internalServerError
	}

	defer rows.Close()

	items := []entity.Item{}

	for rows.Next() {
		item := entity.Item{}
		err = rows.Scan(&item.ItemId, &item.ItemCode, &item.Quantity, &item.Description, &item.CreatedAt, &item.UpdatedAt, &item.OrderId)

		if err != nil {
			return nil, internalServerError
		}

		items = append(items, item)
	}

	return items, nil
}

func (o *orderPg) GetAllOrders() ([]order_repository.OrderItem, error) {
	internalServerError := errors.New("something went wrong")

	rows, err := o.db.Query(getAllOrdersQuery)

	if err != nil {
		return nil, internalServerError
	}

	defer rows.Close()

	orders := []order_repository.OrderItem{}

	for rows.Next() {
		order := entity.Order{}
		err = rows.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt, &order.CreatedAt, &order.UpdatedAt)

		if err != nil {
			return nil, internalServerError
		}

		items, err := o.getItemsByOrderId(order.OrderId)

		if err != nil {
			return nil, internalServerError
		}

		orderItem := order_repository.OrderItem{
			Order: order,
			Items: items,
		}

		orders = append(orders, orderItem)
	}

	return orders, nil
}

func (o *orderPg) UpdateOrder(orderPayload entity.Order, itemsPayload []entity.Item) (*order_repository.OrderItem, error) {
	internalServerError := errors.New("something went wrong")

	tx, err := o.db.Begin()

	if err != nil {

		return nil, internalServerError
	}

	row := tx.QueryRow(updateOrderQuery, orderPayload.OrderId, orderPayload.OrderedAt, orderPayload.CustomerName)

	order := entity.Order{}

	err = row.Scan(&order.OrderId, &order.CustomerName, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		tx.Rollback()
		return nil, internalServerError
	}

	items := []entity.Item{}
	for _, eachItem := range itemsPayload {
		row = tx.QueryRow(updateItemQuery, eachItem.ItemCode, eachItem.Description, eachItem.Quantity)
		item := entity.Item{}
		err = row.Scan(&item.ItemId, &item.ItemCode, &item.Quantity, &item.Description, &item.UpdatedAt, &item.OrderId, &item.CreatedAt)

		if err != nil {
			tx.Rollback()
			return nil, internalServerError
		}

		items = append(items, item)
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return nil, internalServerError
	}

	result := order_repository.OrderItem{
		Order: order,
		Items: items,
	}

	return &result, nil
}

func (o *orderPg) CreateOrder(orderPayload entity.Order, itemsPayload []entity.Item) (*entity.Order, error) {
	tx, err := o.db.Begin()

	internalServerError := errors.New("something went wrong")
	if err != nil {
		return nil, internalServerError
	}

	orderRow := tx.QueryRow(createOrderQuery, orderPayload.OrderedAt, orderPayload.CustomerName)

	var order entity.Order

	err = orderRow.Scan(&order.OrderId, &order.CustomerName, &order.OrderedAt, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		tx.Rollback()
		return nil, internalServerError
	}

	items := []int{}

	for _, eachItem := range itemsPayload {
		itemRow := tx.QueryRow(createItemQuery, eachItem.ItemCode, eachItem.Quantity, eachItem.Description, order.OrderId)

		var itemId int

		err = itemRow.Scan(&itemId)

		if err != nil {
			tx.Rollback()
			return nil, internalServerError
		}

		items = append(items, itemId)
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return nil, internalServerError
	}

	return &order, nil

}

//DeleteOrder return 2 values, first is the number of rows affected, second is the error
func (o *orderPg) DeleteOrder(orderId int) (*http.Response, error) {
	internalServerError := errors.New("something went wrong")

	tx, err := o.db.Begin()

	if err != nil {
		return nil, internalServerError
	}

	result, err := tx.Exec(deleteOrderQuery, orderId)

	if err != nil {
		tx.Rollback()
		return nil, internalServerError
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		tx.Rollback()
		return nil, internalServerError
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("order not found")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return nil, internalServerError
	}

	return &http.Response{
		StatusCode: http.StatusOK,
	}, nil
}