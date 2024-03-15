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

type ItemsInteractor interface {
	CreateItem(item *api.Item) (*repository.ItemModel, error)
	RemoveItem(item *api.Item) (*repository.ItemModel, error)
	UpdateItem(item *api.Item) (*repository.ItemModel, error)
	GetList() ([]*repository.ItemModel, error)
}

type itemsInteractor struct {
	db              *postgres.DB
	itemsRepository repository.ItemsRepository
	pubSub          queue.PubSub
	redis           *redis.Client
	logger          logger.Logger
}

const itemCache = "itemCache"

func NewItemsInteractor(itemsRepository repository.ItemsRepository, redis *redis.Client, pubSub queue.PubSub, logger logger.Logger) ItemsInteractor {
	return &itemsInteractor{
		itemsRepository: itemsRepository,
		redis:           redis,
		pubSub:          pubSub,
		logger:          logger,
	}
}

func (i *itemsInteractor) CreateItem(item *api.Item) (*repository.ItemModel, error) {
	itemDTO := repository.NewItemCreateModel(item.CampaignId, item.Name)
	itemModel, err := i.itemsRepository.Create(itemDTO)

	data, _ := json.Marshal(itemModel)
	_ = i.pubSub.Pub(natsClient.EventTopicName, data)

	return itemModel, err
}

func (i *itemsInteractor) GetList() ([]*repository.ItemModel, error) {
	cacheBytes, err := i.redis.Get(itemCache).Bytes()
	if errors.Is(err, redis.Nil) {
		items, err := i.itemsRepository.GetList()
		if err == nil {
			itemsBytes, _ := json.Marshal(items)
			i.redis.Set(itemCache, interface{}(itemsBytes), time.Minute).Err()
		}

		return items, err
	}

	var items []*repository.ItemModel
	err = json.Unmarshal(cacheBytes, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (i *itemsInteractor) RemoveItem(item *api.Item) (*repository.ItemModel, error) {
	itemDTO := repository.NewItemRemoveModel(item.Id, item.CampaignId)

	itemModel, err := i.itemsRepository.Remove(itemDTO)
	data, _ := json.Marshal(itemModel)
	_ = i.pubSub.Pub(natsClient.EventTopicName, data)

	return itemModel, err
}

func (i *itemsInteractor) UpdateItem(item *api.Item) (*repository.ItemModel, error) {
	itemDTO := repository.NewItemUpdateModel(item.Id, item.CampaignId, item.Name, item.Description)

	itemModel, err := i.itemsRepository.Update(itemDTO)
	data, _ := json.Marshal(itemModel)
	_ = i.pubSub.Pub(natsClient.EventTopicName, data)

	return itemModel, err
}
