package repository

import "time"

type EventsModel struct {
	Id          int
	CampaignId  int
	Name        string
	Description string
	Priority    int
	Removed     bool
	EventTime   time.Time
}

func ItemModelToEvent(itemModel ItemModel) *EventsModel {
	return &EventsModel{
		Id:          itemModel.Id,
		CampaignId:  itemModel.CampaignId,
		Name:        itemModel.Name,
		Description: itemModel.Description,
		Priority:    itemModel.Priority,
		Removed:     itemModel.Removed,
		EventTime:   time.Now(),
	}
}

type EventsRepository interface {
	Create(eventModel *EventsModel) error
}
