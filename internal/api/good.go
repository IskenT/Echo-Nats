package api

import (
	"rest_clickhouse/internal/infrastructure/usecase/repository"
	"time"
)

type Good struct {
	Id          int        `json:"id,omitempty"`
	ProjectId   int        `json:"projectId,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Priority    int        `json:"priority,omitempty"`
	Removed     bool       `json:"removed,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

type GoodList struct {
	Meta Meta `json:"meta"`

	Goods []Good `json:"goods"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

func GetGoodList(GoodModels repository.GoodModelList) GoodList {
	goodList := GoodList{
		Meta: Meta{
			Total:   GoodModels.Meta.Total,
			Removed: GoodModels.Meta.Removed,
			Limit:   GoodModels.Meta.Limit,
			Offset:  GoodModels.Meta.Offset,
		},
		Goods: make([]Good, len(GoodModels.Goods)),
	}

	for i, GoodModel := range GoodModels.Goods {
		good := Good{
			Id:          GoodModel.Id,
			ProjectId:   GoodModel.ProjectId,
			Name:        GoodModel.Name,
			Description: GoodModel.Description,
			Priority:    GoodModel.Priority,
			Removed:     GoodModel.Removed,
			CreatedAt:   &GoodModel.CreatedAt,
		}
		goodList.Goods[i] = good
	}

	return goodList
}

func GetRemovedGood(GoodModel *repository.GoodModel) Good {
	return Good{
		Id:        GoodModel.Id,
		ProjectId: GoodModel.ProjectId,
		Removed:   GoodModel.Removed,
	}
}

func GetUpdatedGood(GoodModel *repository.GoodModel) Good {
	return Good{
		Id:        GoodModel.Id,
		ProjectId: GoodModel.ProjectId,
		Name:      GoodModel.Name,
		Priority:  GoodModel.Priority,
		Removed:   GoodModel.Removed,
		CreatedAt: &GoodModel.CreatedAt,
	}
}
