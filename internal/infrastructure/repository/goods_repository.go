package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rest_clickhouse/internal/infrastructure/usecase/repository"
	postgres "rest_clickhouse/pkg/db"
	"rest_clickhouse/pkg/logger"
	"sync"

	"github.com/go-redis/redis"
)

var (
	ErrGoodNotExist    = errors.New("good not exist")
	ErrProjectNotExist = errors.New("project not exist")
	ErrOnUpdateGood    = errors.New("error when update good")
)

const redisGoodPostfix = "good"

type GoodsRepository struct {
	db          *postgres.DB
	redisClient *redis.Client
	mu          sync.RWMutex
	logger      logger.Logger
}

func NewGoodsRepository(db *postgres.DB, redisClient *redis.Client, logger logger.Logger) repository.GoodsRepository {
	return &GoodsRepository{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (r *GoodsRepository) Create(good *repository.GoodModel) (*repository.GoodModel, error) {
	r.logger.Info("create good")
	var isGoodExist bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT id FROM projects WHERE id = $1)", good.ProjectId).Scan(&isGoodExist)
	if err != nil {
		return nil, err
	}

	if !isGoodExist {
		return nil, ErrProjectNotExist
	}

	query := "select max(priority) FROM goods"
	var maxPriority int
	err = r.db.QueryRow(query).Scan(&maxPriority)
	if err != nil {
		maxPriority = 0
	}

	query = "INSERT INTO goods (project_id, name, description, priority, removed) values ($1, $2, $3, $4, $5) RETURNING id"
	err = r.db.QueryRow(query, good.ProjectId, good.Name, good.Description, maxPriority+1, good.Removed).Scan(&good.Id)
	if err != nil {
		return nil, err
	}

	return r.getGoodById(good.Id, good.ProjectId)
}

func (r *GoodsRepository) GetList(limit, offset int) (*repository.GoodModelList, error) {
	r.logger.Info("get goods")
	goodListModels := &repository.GoodModelList{}
	goodModels := make([]*repository.GoodModel, 0)

	goodQuery := "SELECT id, project_id, name, description, priority, removed, created_at FROM goods ORDER BY id OFFSET $1 LIMIT COALESCE(NULLIF($2, 0), 10)"
	rows, err := r.db.Query(goodQuery, offset, limit)
	if err != nil {
		return goodListModels, err
	}

	for rows.Next() {
		goodModel := new(repository.GoodModel)
		err = rows.Scan(
			&goodModel.Id,
			&goodModel.ProjectId,
			&goodModel.Name,
			&goodModel.Description,
			&goodModel.Priority,
			&goodModel.Removed,
			&goodModel.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		goodModels = append(goodModels, goodModel)
	}

	meta := repository.Meta{
		Total:   len(goodModels),
		Removed: countRemovedGoods(goodModels),
		Limit:   limit,
		Offset:  offset,
	}

	goodListModels.Goods = goodModels
	goodListModels.Meta = meta
	return goodListModels, nil
}

func (r *GoodsRepository) Remove(good *repository.GoodModel) (*repository.GoodModel, error) {
	r.logger.Info("remove good")
	var isGoodExist bool
	r.mu.RLock()
	defer r.mu.RUnlock()

	invalidateKey := fmt.Sprintf("%s-%d", redisGoodPostfix, good.Id)
	r.redisClient.Set(invalidateKey, good, 0)

	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			r.logger.ErrorF("rollback error")
		}
	}()

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow("SELECT EXISTS (SELECT id FROM goods WHERE id = $1 AND project_id = $2)", good.Id, good.ProjectId).Scan(&isGoodExist)
	if err != nil {
		return nil, err
	}

	if !isGoodExist {
		return nil, ErrGoodNotExist
	}

	_, err = tx.Exec("UPDATE goods SET removed = $1 WHERE id = $2 AND project_id = $3", good.Removed, good.Id, good.ProjectId)
	if err != nil {
		return nil, ErrOnUpdateGood
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.getGoodById(good.Id, good.ProjectId)
}

func (r *GoodsRepository) Update(good *repository.GoodModel) (*repository.GoodModel, error) {
	r.logger.Info("update good")

	exists, err := r.checkGoodExistence(good.Id, good.ProjectId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrGoodNotExist
	}

	invalidateKey := fmt.Sprintf("%s-%d", redisGoodPostfix, good.Id)
	if err := r.redisClient.Set(invalidateKey, good, 0).Err(); err != nil {
		return nil, err
	}

	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			r.logger.ErrorF("rollback error")
		}
	}()

	if _, err := tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ"); err != nil {
		return nil, err
	}

	var updateQuery string
	if good.Description != "" {
		updateQuery = "UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id = $4"
	} else {
		updateQuery = "UPDATE goods SET name = $1 WHERE id = $2 AND project_id = $3"
	}
	_, err = tx.Exec(updateQuery, good.Name, good.Description, good.Id, good.ProjectId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.getGoodById(good.Id, good.ProjectId)
}

func (r *GoodsRepository) checkGoodExistence(id, projectId int) (bool, error) {
	var exists bool
	if err := r.db.QueryRow("SELECT EXISTS (SELECT id FROM goods WHERE id = $1 AND project_id = $2)", id, projectId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *GoodsRepository) getGoodById(id, projectId int) (*repository.GoodModel, error) {
	goodModel := new(repository.GoodModel)
	err := r.db.QueryRow("SELECT * FROM goods WHERE id = $1 AND project_id = $2", id, projectId).Scan(
		&goodModel.Id,
		&goodModel.ProjectId,
		&goodModel.Name,
		&goodModel.Description,
		&goodModel.Priority,
		&goodModel.Removed,
		&goodModel.CreatedAt,
	)

	return goodModel, err
}

func countRemovedGoods(goods []*repository.GoodModel) int {
	count := 0
	for _, good := range goods {
		if good.Removed {
			count++
		}
	}
	return count
}
