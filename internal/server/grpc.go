package server

import (
	v1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kbk-authorization/internal/conf"
	"github.com/ZQCard/kbk-authorization/internal/service"
	"github.com/ZQCard/kbk-authorization/pkg/middleware/requestInfo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, service *service.AuthorizationService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			validate.Validator(),
			recovery.Recovery(),
			// 链路追踪
			tracing.Server(),
			// 访问日志
			logging.Server(logger),
			// 元信息
			metadata.Server(),
			// 请求信息中间件
			requestInfo.SetRequestInfo(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterAuthorizationServiceServer(srv, service)
	return srv
}
