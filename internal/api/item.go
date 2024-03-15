package api

import (
	"rest_clickhouse/internal/infrastructure/usecase/repository"
	"time"
)

// для реквестов и ресопнсов с контроллера
type Item struct {
	Id          int        `json:"id,omitempty"`
	CampaignId  int        `json:"campaignId,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Priority    int        `json:"priority,omitempty"`
	Removed     bool       `json:"removed,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

func GetItemList(ItemModels []*repository.ItemModel) []Item {
	itemList := make([]Item, 0)
	for _, ItemModel := range ItemModels {
		item := Item{
			Id:          ItemModel.Id,
			CampaignId:  ItemModel.CampaignId,
			Name:        ItemModel.Name,
			Description: ItemModel.Description,
			Priority:    ItemModel.Priority,
			Removed:     ItemModel.Removed,
			CreatedAt:   &ItemModel.CreatedAt,
		}
		itemList = append(itemList, item)
	}

	return itemList
}

func GetRemovedItem(ItemModel *repository.ItemModel) Item {
	return Item{
		Id:         ItemModel.Id,
		CampaignId: ItemModel.CampaignId,
		Removed:    ItemModel.Removed,
	}
}

func GetUpdatedItem(ItemModel *repository.ItemModel) Item {
	return Item{
		Id:         ItemModel.Id,
		CampaignId: ItemModel.CampaignId,
		Name:       ItemModel.Name,
		Priority:   ItemModel.Priority,
		Removed:    ItemModel.Removed,
		CreatedAt:  &ItemModel.CreatedAt,
	}
}
