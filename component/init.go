package component

import (
	"sync"

	"github.com/xhyonline/xutil/kv"
)

// Server 组件服务
type Server struct {
	Redis *kv.RClient
}

var (
	Instance *Server
	once     sync.Once
)

type Option func()

// Init 初始化组建服务
func Init(options ...Option) {
	once.Do(func() {
		Instance = new(Server)
		for _, f := range options {
			f()
		}
	})
}
