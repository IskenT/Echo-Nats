package repository

import "time"

type EventsModel struct {
	Id          int
	ProjectId   int
	Name        string
	Description string
	Priority    int
	Removed     bool
	EventTime   time.Time
}

func GoodModelToEvent(goodModel GoodModel) *EventsModel {
	return &EventsModel{
		Id:          goodModel.Id,
		ProjectId:   goodModel.ProjectId,
		Name:        goodModel.Name,
		Description: goodModel.Description,
		Priority:    goodModel.Priority,
		Removed:     goodModel.Removed,
		EventTime:   time.Now(),
	}
}

type EventsRepository interface {
	Create(eventModel *EventsModel) error
}
