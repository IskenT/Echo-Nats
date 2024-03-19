package interactors

import (
	"encoding/json"
	"errors"
	"rest_clickhouse/internal/api"
	"rest_clickhouse/internal/infrastructure/queue"
	natsClient "rest_clickhouse/internal/infrastructure/queue/nats"
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
	goodDTO := repository.NewGoodCreateModel(good.ProjectId, good.Name)
	goodModel, err := i.goodsRepository.Create(goodDTO)

	data, _ := json.Marshal(goodModel)
	_ = i.pubSub.Pub(natsClient.EventTopicName, data)

	return goodModel, err
}

func (i *goodsInteractor) GetList(limit, offset int) (*repository.GoodModelList, error) {
	cacheBytes, err := i.redis.Get(goodCache).Bytes()
	if errors.Is(err, redis.Nil) {
		goods, err := i.goodsRepository.GetList(limit, offset)
		if err == nil {
			goodsBytes, _ := json.Marshal(goods)
			if setErr := i.redis.Set(goodCache, interface{}(goodsBytes), time.Minute).Err(); setErr != nil {
				return nil, setErr
			}
		}
		return goods, err
	}

	var goods *repository.GoodModelList
	err = json.Unmarshal(cacheBytes, &goods)
	if err != nil {
		return nil, err
	}

	return goods, nil
}

func (i *goodsInteractor) RemoveGood(good *api.Good) (*repository.GoodModel, error) {
	goodDTO := repository.NewGoodRemoveModel(good.Id, good.ProjectId)

	goodModel, err := i.goodsRepository.Remove(goodDTO)
	data, _ := json.Marshal(goodModel)
	_ = i.pubSub.Pub(natsClient.EventTopicName, data)

	return goodModel, err
}

func (i *goodsInteractor) UpdateGood(good *api.Good) (*repository.GoodModel, error) {
	goodDTO := repository.NewGoodUpdateModel(good.Id, good.ProjectId, good.Name, good.Description)

	goodModel, err := i.goodsRepository.Update(goodDTO)
	data, _ := json.Marshal(goodModel)
	_ = i.pubSub.Pub(natsClient.EventTopicName, data)

	return goodModel, err
}
