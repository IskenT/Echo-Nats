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

func NewGoodsRepository(ctx context.Context, db *postgres.DB, redisClient *redis.Client, logger logger.Logger) repository.GoodsRepository {
	return &GoodsRepository{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (r *GoodsRepository) Create(ctx context.Context, good *repository.GoodModel) (*repository.GoodModel, error) {
	r.logger.Info("create good")
	q := "INSERT INTO goods (project_id, name, description, priority, removed) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var id int
	err := r.db.QueryRow(ctx, q, good.ProjectId, good.Name, good.Description, good.Priority, good.Removed).Scan(&id)
	if err != nil {
		r.logger.ErrorF("error on create good: %v", err)
		return nil, fmt.Errorf("error on create good: %w", err)
	}

	createdGood, err := r.getGoodById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting created good: %w", err)
	}

	return createdGood, nil
}

func (r *GoodsRepository) GetList(ctx context.Context, limit, offset int) (*repository.GoodModelList, error) {
	r.logger.Info("get goods")

	goodListModels := &repository.GoodModelList{}
	goodModels := make([]*repository.GoodModel, 0)

	goodsQuery := "SELECT id, project_id, name, description, priority, removed, created_at FROM goods ORDER BY id OFFSET $1 LIMIT COALESCE(NULLIF($2, 0), 10)"

	rows, err := r.db.Query(ctx, goodsQuery, offset, limit)

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
			return nil, fmt.Errorf("error scanning results: %w", err)
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

func (r *GoodsRepository) Remove(ctx context.Context, good *repository.GoodModel) (*repository.GoodModel, error) {
	r.logger.Info("remove good")

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			r.logger.ErrorF("rollback error")
		}
	}()

	if _, err := tx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ"); err != nil {
		return nil, fmt.Errorf("error setting isolation level: %w", err)
	}

	var isGoodExist bool
	err = tx.QueryRow(ctx, "SELECT EXISTS (SELECT id FROM goods WHERE id = $1 AND project_id = $2)", good.Id, good.ProjectId).Scan(&isGoodExist)
	if err != nil {
		return nil, fmt.Errorf("error checking good existence: %w", err)
	}

	if !isGoodExist {
		return nil, ErrGoodNotExist
	}

	_, err = tx.Exec(ctx, "UPDATE goods SET removed = $1 WHERE id = $2 AND project_id = $3", good.Removed, good.Id, good.ProjectId)
	if err != nil {
		return nil, ErrOnUpdateGood
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	updatedGood, err := r.getGoodById(ctx, good.Id)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated good: %w", err)
	}

	return updatedGood, nil
}

func (r *GoodsRepository) Update(ctx context.Context, good *repository.GoodModel) (*repository.GoodModel, error) {
	r.logger.Info("update good")

	exists, err := r.checkGoodExistence(ctx, good.Id, good.ProjectId)
	if err != nil {
		return nil, fmt.Errorf("error good with same Id exsists: %w", err)
	}
	if !exists {
		return nil, ErrGoodNotExist
	}

	invalidateKey := fmt.Sprintf("%s-%d", redisGoodPostfix, good.Id)
	if err := r.redisClient.Set(invalidateKey, good, 0).Err(); err != nil {
		return nil, fmt.Errorf("error invalidating key: %w", err)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			r.logger.ErrorF("rollback error")
		}
	}()

	if _, err := tx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL REPEATABLE READ"); err != nil {
		return nil, err
	}

	var updateQuery string
	if good.Description != "" {
		updateQuery = "UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id = $4"
	} else {
		updateQuery = "UPDATE goods SET name = $1 WHERE id = $2 AND project_id = $3"
	}
	_, err = tx.Exec(ctx, updateQuery, good.Name, good.Description, good.Id, good.ProjectId)
	if err != nil {
		return nil, fmt.Errorf("error on update: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error on commit: %w", err)
	}

	return r.getGoodById(ctx, good.Id)
}

func (r *GoodsRepository) getGoodById(ctx context.Context, id int) (*repository.GoodModel, error) {
	goodModel := &repository.GoodModel{}

	rows, err := r.db.Query(ctx, "SELECT * FROM goods WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&goodModel.Id, &goodModel.ProjectId, &goodModel.Name, &goodModel.Description, &goodModel.Priority, &goodModel.Removed, &goodModel.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning results: %w", err)
		}
	} else {
		return nil, fmt.Errorf("error scanning results: %w", err)
	}

	return goodModel, nil
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

func (r *GoodsRepository) checkGoodExistence(ctx context.Context, id, projectId int) (bool, error) {
	var exists bool
	if err := r.db.QueryRow(ctx, "SELECT EXISTS (SELECT id FROM goods WHERE id = $1 AND project_id = $2)", id, projectId).Scan(&exists); err != nil {
		return false, fmt.Errorf("error on exsistance: %w", err)
	}
	return exists, nil
}
