package http

import (
	"context"
	"fmt"
	httpControllers "rest_clickhouse/internal/infrastructure/interfaces"
	"rest_clickhouse/pkg/logger"

	"github.com/labstack/echo/v4"
)

type HTTPServer interface {
	Start()
	Stop(ctx context.Context)
}

type EchoHTTPServer struct {
	echo            *echo.Echo
	serverPort      string
	itemsController httpControllers.ItemsController
	logger          logger.Logger
}

func NewEchoHTTPServer(
	ServerPort string,
	itemsController httpControllers.ItemsController,
	logger logger.Logger,
) *EchoHTTPServer {
	server := &EchoHTTPServer{
		echo:            echo.New(),
		itemsController: itemsController,
		serverPort:      ServerPort,
		logger:          logger,
	}

	return server
}

func (s *EchoHTTPServer) Start() {

	s.echo.POST("/items/create/:campaignId", s.handleCreateItem)
	s.echo.GET("/items/list", s.handleGetItems)
	s.echo.DELETE("/item/remove/:id/:campaignId", s.handleRemoveItem)
	s.echo.PATCH("/item/update/:id/:campaignId", s.handleUpdateItem)

	func() {
		port := fmt.Sprintf(":%v", s.serverPort)
		if err := s.echo.Start(port); err != nil {
			s.logger.Error("Echo error:", err)
		}
	}()
}

func (s *EchoHTTPServer) Stop(ctx context.Context) {
	err := s.echo.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Echo error:", err)
	}
}

func (s *EchoHTTPServer) handleCreateItem(ctx echo.Context) error {
	return s.itemsController.HandleCreateItem(ctx)
}

func (s *EchoHTTPServer) handleGetItems(ctx echo.Context) error {
	return s.itemsController.HandleGetItem(ctx)
}

func (s *EchoHTTPServer) handleRemoveItem(ctx echo.Context) error {
	return s.itemsController.HandleRemoveItem(ctx)
}

func (s *EchoHTTPServer) handleUpdateItem(ctx echo.Context) error {
	return s.itemsController.HandleUpdateItems(ctx)
}
