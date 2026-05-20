//go:build wireinject
// +build wireinject

package main

import (
	"admin/internal/biz"
	"admin/internal/conf"
	"admin/internal/data"
	"admin/internal/server"
	"admin/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp Wire 注入器声明
func wireApp(bc *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	wire.Build(
		// 从 Bootstrap 提取子配置
		provideServerConf,
		provideDataConf,
		provideAuthConf,

		// 四层注入
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,

		// App
		newApp,
	)
	return nil, nil, nil
}
