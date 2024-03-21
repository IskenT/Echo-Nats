package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rest_clickhouse/internal/infrastructure/usecase/repository"
	"rest_clickhouse/pkg/logger"
)

const eventsPackCount = 100

type EventsRepository struct {
	clickHouseConn *sql.DB
	eventModels    []*repository.EventsModel
	logger         logger.Logger
}

func NewLogsRepository(clickHouseConn *sql.DB, logger logger.Logger) repository.EventsRepository {
	return &EventsRepository{
		clickHouseConn: clickHouseConn,
		eventModels:    make([]*repository.EventsModel, 0),
		logger:         logger,
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
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			r.logger.ErrorF("rollback error: %v", err)
		}
	}()

	query := "INSERT INTO Events (Id,ProjectId,Name,Description,Priority,Removed,EventTime) values ($1, $2,$3,$4,$5,$6,$7)"
	for _, event := range r.eventModels {
		_, err = tx.Exec(
			query,
			event.Id,
			event.ProjectId,
			event.Name,
			event.Description,
			event.Priority,
			event.Removed,
			event.EventTime)
		if err != nil {
			return fmt.Errorf("error executing query: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	// Очищаем список eventModels после успешного коммита.
	r.eventModels = r.eventModels[:0]

	return nil
}
