package server

import "github.com/google/wire"

// ProviderSet server 层依赖注入集合
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)
