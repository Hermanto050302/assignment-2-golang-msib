package item_repository

import "assignment-2-golang-msib/entity"

type ItemRepository interface {
	FindItemsByItemCodes(itemCodes []string) ([]*entity.Item, error)
}
