package repository

import (
	"context"
	"time"
)

// GoodModel содержит информацию о товаре.
type GoodModel struct {
	Id          int       `db:"id"`
	ProjectId   int       `db:"project_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	CreatedAt   time.Time `db:"created_at"`
}

// GoodModelList содержит список товаров и метаданные.
type GoodModelList struct {
	Meta  Meta
	Goods []*GoodModel
}

// Meta содержит метаданные.
type Meta struct {
	Total   int // Общее количество элементов
	Removed int // Количество удаленных элементов
	Limit   int // Ограничение количества элементов в списке
	Offset  int // Смещение для пагинации
}

func NewGoodCreateModel(projectId int, name string) *GoodModel {
	return &GoodModel{
		ProjectId:   projectId,
		Name:        name,
		Description: "",
		Removed:     false,
	}
}

func NewGoodUpdateModel(id int, projectId int, name string, description string) *GoodModel {
	return &GoodModel{
		Id:          id,
		ProjectId:   projectId,
		Name:        name,
		Description: description,
	}
}

func NewGoodRemoveModel(id int, projectId int) *GoodModel {
	return &GoodModel{
		Id:        id,
		ProjectId: projectId,
		Removed:   true,
	}
}

type GoodsRepository interface {
	Create(ctx context.Context, Good *GoodModel) (*GoodModel, error)
	GetList(ctx context.Context, limit, offset int) (*GoodModelList, error)
	Remove(ctx context.Context, good *GoodModel) (*GoodModel, error)
	Update(ctx context.Context, good *GoodModel) (*GoodModel, error)
}
