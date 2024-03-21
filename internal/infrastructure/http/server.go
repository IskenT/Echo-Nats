package http

import (
	"context"
	"fmt"
	"rest_clickhouse/pkg/logger"

	"github.com/labstack/echo/v4"
)

type HTTPServer interface {
	Start()
	Stop(ctx context.Context)
}

type EchoHTTPServer struct {
	echo         *echo.Echo
	serverPort   string
	goodsService GoodsService
	logger       logger.Logger
}

func NewEchoHTTPServer(
	ServerPort string,
	goodsService GoodsService,
	logger logger.Logger,
) *EchoHTTPServer {
	server := &EchoHTTPServer{
		echo:         echo.New(),
		goodsService: goodsService,
		serverPort:   ServerPort,
		logger:       logger,
	}

	return server
}

func (s *EchoHTTPServer) Start() {
	s.echo.POST("/goods/create/:projectId", s.handleCreateGood)
	s.echo.GET("/goods/list/:limit/:offset", s.handleGetGoods)
	s.echo.DELETE("/good/remove/:id/:projectId", s.handleRemoveGood)
	s.echo.PATCH("/good/update/:id/:projectId", s.handleUpdateGood)

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

func (s *EchoHTTPServer) handleCreateGood(ctx echo.Context) error {
	return s.goodsService.HandleCreateGood(ctx)
}

func (s *EchoHTTPServer) handleGetGoods(ctx echo.Context) error {
	return s.goodsService.HandleGetGood(ctx)
}

func (s *EchoHTTPServer) handleRemoveGood(ctx echo.Context) error {
	return s.goodsService.HandleRemoveGood(ctx)
}

func (s *EchoHTTPServer) handleUpdateGood(ctx echo.Context) error {
	return s.goodsService.HandleUpdateGoods(ctx)
}
