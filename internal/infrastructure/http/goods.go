package http

import (
	"errors"
	"fmt"
	"net/http"
	"rest_clickhouse/internal/api"
	repository2 "rest_clickhouse/internal/infrastructure/repository"
	"rest_clickhouse/internal/infrastructure/usecase/interactors"
	"rest_clickhouse/pkg/logger"
	"strconv"

	"github.com/labstack/echo/v4"
)

type GoodsService interface {
	HandleCreateGood(c echo.Context) error
	HandleGetGood(ctx echo.Context) error
	HandleRemoveGood(ctx echo.Context) error
	HandleUpdateGoods(ctx echo.Context) error
}

type goodsService struct {
	goodsInteractor interactors.GoodsInteractor
	logger          logger.Logger
}

func NewGoodsService(goodsInteractor interactors.GoodsInteractor, logger logger.Logger) GoodsService {
	return &goodsService{
		goodsInteractor: goodsInteractor,
		logger:          logger,
	}
}

func (c *goodsService) HandleCreateGood(ctx echo.Context) error {
	good := new(api.Good)
	projectId, err := strconv.Atoi(ctx.Param("projectId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	err = ctx.Bind(&good)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid body")
	}

	good.ProjectId = projectId

	goodDTO, err := c.goodsInteractor.CreateGood(good)
	if errors.Is(err, repository2.ErrProjectNotExist) {
		return ctx.String(http.StatusNotFound, "ProjectId not found")
	}

	if err != nil {
		fmt.Printf("error on create good: %w", err)
		return ctx.String(http.StatusInternalServerError, "internal error")
	}

	response := api.GetUpdatedGood(goodDTO)
	return ctx.JSON(http.StatusCreated, response)
}

func (c *goodsService) HandleGetGood(ctx echo.Context) error {
	limit, err := strconv.Atoi(ctx.Param("limit"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid limit")
	}

	offset, err := strconv.Atoi(ctx.Param("offset"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid offset")
	}

	goodsModelList, err := c.goodsInteractor.GetList(limit, offset)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "internal error")
	}

	goodsList := api.GetGoodList(*goodsModelList)

	return ctx.JSON(http.StatusOK, goodsList)
}

func (c *goodsService) HandleRemoveGood(ctx echo.Context) error {
	good := new(api.Good)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	projectId, err := strconv.Atoi(ctx.Param("projectId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	good.Id = id
	good.ProjectId = projectId

	goodDTO, err := c.goodsInteractor.RemoveGood(good)
	if errors.Is(err, repository2.ErrGoodNotExist) {
		return ctx.JSON(http.StatusNotFound, api.NewErrorResponse(api.GoodNotFoundCode, api.GoodNotFoundMessage))
	}

	if err != nil {
		return ctx.String(http.StatusInternalServerError, "internal error")
	}

	response := api.GetRemovedGood(goodDTO)

	return ctx.JSON(http.StatusOK, response)
}

func (c *goodsService) HandleUpdateGoods(ctx echo.Context) error {
	good := new(api.Good)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	projectId, err := strconv.Atoi(ctx.Param("projectId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	err = ctx.Bind(&good)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid body")
	}

	if good.Name == "" {
		return ctx.String(http.StatusBadRequest, "invalid name")
	}

	good.Id = id
	good.ProjectId = projectId

	goodDTO, err := c.goodsInteractor.UpdateGood(good)
	if errors.Is(err, repository2.ErrGoodNotExist) {
		return ctx.JSON(http.StatusNotFound, api.NewErrorResponse(api.GoodNotFoundCode, api.GoodNotFoundMessage))
	}

	if err != nil {
		return ctx.String(http.StatusInternalServerError, "internal error")
	}

	response := api.GetUpdatedGood(goodDTO)
	return ctx.JSON(http.StatusOK, response)
}
