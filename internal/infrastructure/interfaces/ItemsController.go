package httpControllers

import (
	"errors"
	"net/http"
	"rest_clickhouse/internal/api"
	repository2 "rest_clickhouse/internal/infrastructure/repository"
	interactors "rest_clickhouse/internal/infrastructure/usecase/interractors"
	"rest_clickhouse/pkg/logger"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ItemsController interface {
	HandleCreateItem(c echo.Context) error
	HandleGetItem(ctx echo.Context) error
	HandleRemoveItem(ctx echo.Context) error
	HandleUpdateItems(ctx echo.Context) error
}

type itemsController struct {
	itemsInteractor interactors.ItemsInteractor
	logger          logger.Logger
}

func NewItemsController(itemsInteractor interactors.ItemsInteractor, logger logger.Logger) ItemsController {
	return &itemsController{
		itemsInteractor: itemsInteractor,
		logger:          logger,
	}
}

func (c *itemsController) HandleCreateItem(ctx echo.Context) error {
	item := new(api.Item)

	campaignId, err := strconv.Atoi(ctx.Param("campaignId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	err = ctx.Bind(&item)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid body")
	}

	item.CampaignId = campaignId

	itemDTO, err := c.itemsInteractor.CreateItem(item)
	if errors.Is(err, repository2.CampaignNotExistError) {
		return ctx.String(http.StatusNotFound, "campaign not found")
	}

	if err != nil {
		return ctx.String(http.StatusInternalServerError, "internal error")
	}

	response := api.GetUpdatedItem(itemDTO)
	return ctx.JSON(http.StatusCreated, response)
}

func (c *itemsController) HandleGetItem(ctx echo.Context) error {
	items, err := c.itemsInteractor.GetList()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "internal error")
	}

	response := api.GetItemList(items)

	return ctx.JSON(http.StatusOK, response)
}

func (c *itemsController) HandleRemoveItem(ctx echo.Context) error {
	item := new(api.Item)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	campaignId, err := strconv.Atoi(ctx.Param("campaignId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	item.Id = id
	item.CampaignId = campaignId

	itemDTO, err := c.itemsInteractor.RemoveItem(item)
	if errors.Is(err, repository2.ItemNotExistError) {
		return ctx.JSON(http.StatusNotFound, api.NewErrorResponse(api.ItemNotFoundCode, api.ItemNotFoundMessage))
	}

	if err != nil {
		return ctx.String(http.StatusInternalServerError, "internal error")
	}

	response := api.GetRemovedItem(itemDTO)

	return ctx.JSON(http.StatusOK, response)
}

func (c *itemsController) HandleUpdateItems(ctx echo.Context) error {
	item := new(api.Item)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	campaignId, err := strconv.Atoi(ctx.Param("campaignId"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid url params")
	}

	err = ctx.Bind(&item)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "invalid body")
	}

	if item.Name == "" {
		return ctx.String(http.StatusBadRequest, "invalid name")
	}

	item.Id = id
	item.CampaignId = campaignId

	itemDTO, err := c.itemsInteractor.UpdateItem(item)
	if errors.Is(err, repository2.ItemNotExistError) {
		return ctx.JSON(http.StatusNotFound, api.NewErrorResponse(api.ItemNotFoundCode, api.ItemNotFoundMessage))
	}

	if err != nil {
		return ctx.String(http.StatusInternalServerError, "internal error")
	}

	response := api.GetUpdatedItem(itemDTO)
	return ctx.JSON(http.StatusOK, response)
}
