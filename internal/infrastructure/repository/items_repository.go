package repository

import (
	"context"
	"errors"
	"fmt"
	"rest_clickhouse/internal/infrastructure/usecase/repository"
	postgres "rest_clickhouse/pkg/db"
	"rest_clickhouse/pkg/logger"
	"sync"

	"github.com/go-redis/redis"
)

type ItemsRepository struct {
	db          *postgres.DB
	redisClient *redis.Client
	mu          sync.RWMutex
	logger      logger.Logger
}

var ItemNotExistError = errors.New("item not exist")
var CampaignNotExistError = errors.New("item not exist")

const redisItemPostfix = "item"

func NewItemsRepository(db *postgres.DB, redisClient *redis.Client, logger logger.Logger) repository.ItemsRepository {
	return &ItemsRepository{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (r *ItemsRepository) Create(item *repository.ItemModel) (*repository.ItemModel, error) {
	r.logger.Info("create item")
	var isItemExist bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT id FROM campaigns WHERE id = $1)", item.CampaignId).Scan(&isItemExist)
	if err != nil {
		return nil, err
	}

	if !isItemExist {
		return nil, CampaignNotExistError
	}

	query := "select max(priority) FROM items"
	var maxPriority int
	err = r.db.QueryRow(query).Scan(&maxPriority)
	if err != nil {
		maxPriority = 0
	}

	query = "INSERT INTO items (campaign_id,name,priority,removed) values ($1, $2,$3,$4) RETURNING id"
	err = r.db.QueryRow(query, item.CampaignId, item.Name, maxPriority+1, item.Removed).Scan(&item.Id)
	if err != nil {
		return nil, err
	}

	return r.getItemById(item.Id, item.CampaignId)
}

func (r *ItemsRepository) GetList() ([]*repository.ItemModel, error) {
	r.logger.Info("get items")
	itemModels := make([]*repository.ItemModel, 0)

	rows, err := r.db.Query("SELECT * FROM items")
	if err != nil {
		r.logger.ErrorF("repo error:%w", err)
	}

	for rows.Next() {
		itemModel := new(repository.ItemModel)
		err = rows.Scan(
			&itemModel.Id,
			&itemModel.CampaignId,
			&itemModel.Name,
			&itemModel.Description,
			&itemModel.Priority,
			&itemModel.Removed,
			&itemModel.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		itemModels = append(itemModels, itemModel)
	}

	return itemModels, err
}

func (r *ItemsRepository) Remove(item *repository.ItemModel) (*repository.ItemModel, error) {
	r.logger.Info("remove item")
	var isItemExist bool
	r.mu.RLock()
	defer r.mu.RUnlock()

	invalidateKey := fmt.Sprintf("%s-%s", redisItemPostfix, string(item.Id))
	r.redisClient.Set(invalidateKey, item, 0)

	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ ")

	err = tx.QueryRow("SELECT EXISTS (SELECT id FROM items WHERE id = $1 AND campaign_id = $2)", item.Id, item.CampaignId).Scan(&isItemExist)
	if err != nil {
		return nil, err
	}

	if !isItemExist {
		return nil, ItemNotExistError
	}

	_, err = tx.Exec("UPDATE items SET removed = $1 WHERE id = $2 AND campaign_id = $3", item.Removed, item.Id, item.CampaignId)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.getItemById(item.Id, item.CampaignId)
}

func (r *ItemsRepository) Update(item *repository.ItemModel) (*repository.ItemModel, error) {
	r.logger.Info("update item")

	var isItemExist bool
	r.mu.RLock()
	defer r.mu.RUnlock()

	invalidateKey := fmt.Sprintf("%s-%s", redisItemPostfix, string(item.Id))
	r.redisClient.Set(invalidateKey, item, 0)

	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ ")
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow("SELECT EXISTS (SELECT id FROM items WHERE id = $1 AND campaign_id = $2)", item.Id, item.CampaignId).Scan(&isItemExist)
	if err != nil {
		return nil, err
	}

	if !isItemExist {
		return nil, ItemNotExistError
	}

	if item.Description != "" {
		_, err = tx.Exec("UPDATE items SET name = $1, description = $2 WHERE id = $3 AND campaign_id = $4", item.Name, item.Description, item.Id, item.CampaignId)
	} else {
		_, err = tx.Exec("UPDATE items SET name = $1 WHERE id = $2 AND campaign_id = $3", item.Name, item.Id, item.CampaignId)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.getItemById(item.Id, item.CampaignId)
}

func (r *ItemsRepository) getItemById(id, campaignId int) (*repository.ItemModel, error) {
	itemModel := new(repository.ItemModel)
	err := r.db.QueryRow("SELECT * FROM items WHERE id = $1 AND campaign_id = $2", id, campaignId).Scan(
		&itemModel.Id,
		&itemModel.CampaignId,
		&itemModel.Name,
		&itemModel.Description,
		&itemModel.Priority,
		&itemModel.Removed,
		&itemModel.CreatedAt,
	)

	return itemModel, err
}
