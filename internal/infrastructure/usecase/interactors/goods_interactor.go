package interactors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"rest_clickhouse/internal/api"
	"rest_clickhouse/internal/infrastructure/queue"
	nats_client "rest_clickhouse/internal/infrastructure/queue/nats"
	"rest_clickhouse/internal/infrastructure/usecase/repository"
	postgres "rest_clickhouse/pkg/db"
	"rest_clickhouse/pkg/logger"
	"time"

	"github.com/go-redis/redis"
)

type GoodsInteractor interface {
	CreateGood(good *api.Good) (*repository.GoodModel, error)
	RemoveGood(good *api.Good) (*repository.GoodModel, error)
	UpdateGood(good *api.Good) (*repository.GoodModel, error)
	GetList(limit, offset int) (*repository.GoodModelList, error)
}

type goodsInteractor struct {
	db              *postgres.DB
	goodsRepository repository.GoodsRepository
	pubSub          queue.PubSub
	redis           *redis.Client
	logger          logger.Logger
}

const goodCache = "goodCache"

func NewGoodsInteractor(goodsRepository repository.GoodsRepository, redis *redis.Client, pubSub queue.PubSub, logger logger.Logger) GoodsInteractor {
	return &goodsInteractor{
		goodsRepository: goodsRepository,
		redis:           redis,
		pubSub:          pubSub,
		logger:          logger,
	}
}

func (i *goodsInteractor) CreateGood(good *api.Good) (*repository.GoodModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	goodDTO := repository.NewGoodCreateModel(good.ProjectId, good.Name)
	goodModel, err := i.goodsRepository.Create(ctx, goodDTO)
	if err != nil {
		return nil, fmt.Errorf("error on create good: %w", err)
	}

	data, err := json.Marshal(goodModel)
	if err != nil {
		return goodModel, fmt.Errorf("error marshaling goodModel: %w", err)
	}

	if err := i.pubSub.Pub(nats_client.EventTopicName, data); err != nil {
		return goodModel, fmt.Errorf("error publishing event: %w", err)
	}

	return goodModel, nil
}

func (i *goodsInteractor) GetList(limit, offset int) (*repository.GoodModelList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cacheBytes, err := i.redis.Get(goodCache).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("error getting data from cache: %w", err)
	}

	if errors.Is(err, redis.Nil) {
		goods, err := i.goodsRepository.GetList(ctx, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("error getting list from repository: %w", err)
		}
		goodsBytes, err := json.Marshal(goods)
		if err != nil {
			return nil, fmt.Errorf("error marshaling goods: %w", err)
		}
		if err := i.redis.Set(goodCache, goodsBytes, time.Minute).Err(); err != nil {
			return nil, fmt.Errorf("error setting data in cache: %w", err)
		}
		return goods, nil
	}

	var goods repository.GoodModelList
	if err := json.Unmarshal(cacheBytes, &goods); err != nil {
		return nil, fmt.Errorf("error unmarshaling cached data: %w", err)
	}

	return &goods, nil
}

func (i *goodsInteractor) RemoveGood(good *api.Good) (*repository.GoodModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	goodDTO := repository.NewGoodRemoveModel(good.Id, good.ProjectId)

	goodModel, err := i.goodsRepository.Remove(ctx, goodDTO)
	if err != nil {
		return nil, fmt.Errorf("error on remove good: %w", err)
	}

	data, err := json.Marshal(goodModel)
	if err != nil {
		return nil, fmt.Errorf("error marshaling goodModel: %w", err)
	}

	if err := i.pubSub.Pub(nats_client.EventTopicName, data); err != nil {
		return nil, fmt.Errorf("error publishing event: %w", err)
	}

	return goodModel, nil
}

func (i *goodsInteractor) UpdateGood(good *api.Good) (*repository.GoodModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	goodDTO := repository.NewGoodUpdateModel(good.Id, good.ProjectId, good.Name, good.Description)

	goodModel, err := i.goodsRepository.Update(ctx, goodDTO)
	if err != nil {
		return nil, fmt.Errorf("error on update good: %w", err)
	}

	data, err := json.Marshal(goodModel)
	if err != nil {
		return nil, fmt.Errorf("error marshaling goodModel: %w", err)
	}

	if err := i.pubSub.Pub(nats_client.EventTopicName, data); err != nil {
		return nil, fmt.Errorf("error publishing event: %w", err)
	}

	return goodModel, nil
}
