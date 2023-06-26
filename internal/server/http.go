package server

import (
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"

	v1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kbk-authorization/internal/conf"
	"github.com/ZQCard/kbk-authorization/internal/service"
	"github.com/ZQCard/kbk-authorization/pkg/middleware/requestInfo"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, service *service.AuthorizationService, tp *tracesdk.TracerProvider, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			validate.Validator(),
			recovery.Recovery(),
			tracing.Server(),
			// 元信息
			metadata.Server(),
			// 访问日志
			logging.Server(logger),
			// 租户信息
			requestInfo.SetRequestInfo(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	openAPIhandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIhandler)
	v1.RegisterAuthorizationServiceHTTPServer(srv, service)
	return srv
}
