package repository

import (
	"time"
)

type GoodModel struct {
	Id          int
	ProjectId   int
	Name        string
	Description string
	Priority    int
	Removed     bool
	CreatedAt   time.Time
}

type GoodModelList struct {
	Meta  Meta
	Goods []*GoodModel
}

type Meta struct {
	Total   int
	Removed int
	Limit   int
	Offset  int
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
	Create(Good *GoodModel) (*GoodModel, error)
	GetList(limit, offset int) (GoodModelList, error)
	Remove(good *GoodModel) (*GoodModel, error)
	Update(good *GoodModel) (*GoodModel, error)
}
