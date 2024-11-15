package configpost

import (
	"context"
	"doh-server/pkg/config"
	"doh-server/server/controller"
	"doh-server/server/service"
)

func PostInit(ctx context.Context) {
	// 初始化客户端
	service.InitClient()
	// 启动 UDP DNS 服务
	if config.Conf.UDPBind != "" {
		go controller.RunDNServer(ctx, config.Conf.UDPBind, "udp")
	}
	// 启动 TCP DNS 服务
	if config.Conf.TCPBind != "" {
		go controller.RunDNServer(ctx, config.Conf.TCPBind, "tcp")
	}
	// 启动 HTTP 服务
	controller.RunServer(ctx, config.Conf.ServerBind)
}
