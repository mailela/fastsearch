package router

import (
	"gitee.com/rachel_os/fastsearch/web/controller"

	"github.com/gin-gonic/gin"
)

// InitNegativeRouter 负面词路由
func InitNegativeRouter(Router *gin.RouterGroup) {

	BaseRouter := Router.Group("negative")
	{
		BaseRouter.POST("add", controller.NegativeKeys)
		BaseRouter.POST("import", controller.NegativeBatch)
		BaseRouter.POST("remove", controller.NegativeRemove)
		BaseRouter.POST("query", controller.QueryNegative)
		BaseRouter.POST("check", controller.HasNegative)
		BaseRouter.POST("keys", controller.AllKeys)
		BaseRouter.POST("apply", controller.AllKeys)

	}
}
