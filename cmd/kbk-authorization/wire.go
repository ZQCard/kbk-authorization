//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/biz"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/conf"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/data"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/server"
	"github.com/ZQCard/kratos-base-kit/kbk-authorization/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// wireApp init kratos application.
func wireApp(*conf.Env, *conf.Server, *conf.Registry, *conf.Data, *conf.Bootstrap, log.Logger, *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
