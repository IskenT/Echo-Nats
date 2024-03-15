package repository

import "time"

type ItemModel struct {
	Id          int
	CampaignId  int
	Name        string
	Description string
	Priority    int
	Removed     bool
	CreatedAt   time.Time
}

func NewItemCreateModel(campaignId int, name string) *ItemModel {
	return &ItemModel{
		CampaignId:  campaignId,
		Name:        name,
		Description: "",
		Removed:     false,
	}
}

func NewItemUpdateModel(id int, campaignId int, name string, description string) *ItemModel {
	return &ItemModel{
		Id:          id,
		CampaignId:  campaignId,
		Name:        name,
		Description: description,
	}
}

func NewItemRemoveModel(id int, campaignId int) *ItemModel {
	return &ItemModel{
		Id:         id,
		CampaignId: campaignId,
		Removed:    true,
	}
}

type ItemsRepository interface {
	Create(Item *ItemModel) (*ItemModel, error)
	GetList() ([]*ItemModel, error)
	Remove(item *ItemModel) (*ItemModel, error)
	Update(item *ItemModel) (*ItemModel, error)
}
