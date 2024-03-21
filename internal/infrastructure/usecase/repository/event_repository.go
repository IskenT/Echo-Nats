package repository

import "time"

type EventsModel struct {
	Id          int    `json:"id" db:"id"`
	ProjectId   int    `json:"projectId" db:"project_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Priority    int    `json:"priority" db:"priority"`
	Removed     bool   `json:"removed" db:"removed"`
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
