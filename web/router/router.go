package router

import (
	"io"
	"log"
	"mime"
	"strings"

	"gitee.com/rachel_os/fastsearch/global"
	"gitee.com/rachel_os/fastsearch/web/admin"
	"gitee.com/rachel_os/fastsearch/web/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// SetupRouter 路由管理
func SetupRouter() *gin.Engine {
	if global.CONFIG.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard //禁止Gin的控制台输出
	}

	router := gin.Default()

	// 启用GZIP压缩
	if global.CONFIG.EnableGzip {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	var handlers []gin.HandlerFunc
	//认证
	if global.CONFIG.Auth != "" {
		auths := strings.Split(global.CONFIG.Auth, ":")
		handlers = append(handlers, gin.BasicAuth(
			gin.Accounts{
				auths[0]: auths[1],
			},
		),
		)
		log.Println("Enable Auth:", global.CONFIG.Auth)
	}
	handlers = append(handlers, func(c *gin.Context) {
		c.Header("Server", "FastSearch")
	})
	// 告诉服务.js文件的MIME类型
	err := mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	// 如果存在错误则需要马上抛出
	if err != nil {
		panic("添加扩展类型 mime 错误，错误原因：" + err.Error())
	}

	//注册admin
	if global.CONFIG.EnableAdmin {
		admin.Register(router, handlers...)
		log.Printf("Admin Url: \t http://%v/admin", global.CONFIG.Addr)
	}

	// 分组管理 中间件管理
	router.Use(middleware.Cors(), middleware.Exception())
	group := router.Group("/api", handlers...)
	{
		InitBaseRouter(group)     // 基础管理
		InitNegativeRouter(group) // 负面管理
		InitIndexRouter(group)    // 索引管理
		InitDatabaseRouter(group) // 数据库管理
		InitWordRouter(group)     // 分词管理
	}
	log.Printf("API Url: \t http://%v/api", global.CONFIG.Addr)
	return router
}
