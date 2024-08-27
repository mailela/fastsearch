package controller

import (
	"gitee.com/rachel_os/fastsearch/searcher/model"

	"github.com/gin-gonic/gin"
)

func Welcome(c *gin.Context) {
	ResponseSuccessWithData(c, "Welcome to fastsearch")
}

// Query 查询
func Query(c *gin.Context) {

	var request = &model.SearchRequest{
		Database: c.Query("database"),
	}
	if err := c.ShouldBind(&request); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	if srv.Database.Container.CheckDataBase(request.Database) == false {
		ResponseErrorWithMsg(c, "数据库不存在")
		return
	}

	//调用搜索
	r, err := srv.Base.Query(request)
	if err != nil {
		ResponseErrorWithMsg(c, err.Error())
	} else {
		ResponseSuccessWithData(c, r)
	}
}

// GC 释放GC
func GC(c *gin.Context) {
	srv.Base.GC()
	ResponseSuccess(c)
}

// Status 获取服务器状态
func Status(c *gin.Context) {
	r := srv.Base.Status()
	ResponseSuccessWithData(c, r)
}
