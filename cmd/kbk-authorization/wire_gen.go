// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/biz"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/conf"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/data"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/server"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/sdk/trace"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(env *conf.Env, confServer *conf.Server, registry *conf.Registry, confData *conf.Data, bootstrap *conf.Bootstrap, logger log.Logger, tracerProvider *trace.TracerProvider) (*kratos.App, func(), error) {
	db := data.NewMysqlCmd(bootstrap, logger)
	client := data.NewRedisClient(confData)
	dataData, cleanup, err := data.NewData(bootstrap, db, client, logger)
	if err != nil {
		return nil, nil, err
	}
	menuRepo := data.NewMenuRepo(dataData, logger)
	menuUsecase := biz.NewMenuUsecase(menuRepo, logger)
	casbinRepo := data.NewCasbinRepo(dataData, bootstrap, logger)
	casbinUsecase := biz.NewCasbinUsecase(casbinRepo, logger)
	roleRepo := data.NewRoleRepo(dataData, logger)
	roleUsecase := biz.NewRoleUsecase(roleRepo, logger)
	apiRepo := data.NewApiRepo(dataData, logger)
	apiUsecase := biz.NewApiUsecase(apiRepo, logger)
	authorizationService := service.NewAuthorizationService(menuUsecase, casbinUsecase, roleUsecase, apiUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, authorizationService, logger)
	httpServer := server.NewHTTPServer(confServer, authorizationService, tracerProvider, logger)
	registrar := data.NewRegistrar(registry)
	app := newApp(logger, grpcServer, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}