package global

import (
	"gitee.com/rachel_os/fastsearch/searcher"
)

const VERSION = "1.0.2"

var (
	CONFIG    *Config // 服务器设置
	Container *searcher.Container
)
