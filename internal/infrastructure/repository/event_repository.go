package repository

import (
	"context"
	"database/sql"
	"rest_clickhouse/internal/infrastructure/usecase/repository"
)

const eventsPackCount = 100

type EventsRepository struct {
	clickHouseConn *sql.DB
	eventModels    []*repository.EventsModel
}

func NewLogsRepository(clickHouseConn *sql.DB) repository.EventsRepository {
	return &EventsRepository{
		clickHouseConn: clickHouseConn,
		eventModels:    make([]*repository.EventsModel, 0),
	}
}

func (r *EventsRepository) Create(eventModel *repository.EventsModel) error {
	if len(r.eventModels) < eventsPackCount {
		r.eventModels = append(r.eventModels, eventModel)
		return nil
	}

	ctx := context.Background()
	tx, err := r.clickHouseConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO Events (Id,CampaignId,Name,Description,Priority,Removed,EventTime) values ($1, $2,$3,$4,$5,$6,$7)"
	for _, event := range r.eventModels {
		_, err = tx.Exec(
			query,
			event.Id,
			event.CampaignId,
			event.Name,
			event.Description,
			event.Priority,
			event.Removed,
			event.EventTime)
		if err != nil {
			//логируем но не возвращаем так как пачка может пойти назад
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	r.eventModels = r.eventModels[:0]

	return err
}
