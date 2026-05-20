package main

import (
	"flag"
	"os"

	"admin/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// 编译时注入
var (
	Name    = "admin"
	Version = "1.0.0"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "cmd/admin/conf/config.yaml", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	)
}

// provideServerConf 提取服务器配置
func provideServerConf(bc *conf.Bootstrap) *conf.Server {
	return &bc.Server
}

// provideDataConf 提取数据源配置
func provideDataConf(bc *conf.Bootstrap) *conf.Data {
	return &bc.Data
}

// provideAuthConf 提取认证配置
func provideAuthConf(bc *conf.Bootstrap) *conf.Auth {
	return &bc.Auth
}

func main() {
	flag.Parse()

	// 初始化日志
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.name", Name,
		"service.version", Version,
	)
	helper := log.NewHelper(logger)

	// 加载配置
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		helper.Fatalf("failed to load config: %v", err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		helper.Fatalf("failed to scan config: %v", err)
	}

	// Wire 注入并启动
	app, cleanup, err := wireApp(&bc, logger)
	if err != nil {
		helper.Fatalf("failed to wire app: %v", err)
	}
	defer cleanup()

	helper.Infof("starting %s %s ...", Name, Version)
	helper.Infof("gRPC server listening on %s", bc.Server.GRPC.Addr)
	helper.Infof("HTTP server listening on %s", bc.Server.HTTP.Addr)

	if err := app.Run(); err != nil {
		helper.Fatalf("failed to run app: %v", err)
	}
}
