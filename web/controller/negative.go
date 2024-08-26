package controller

import (
	"gitee.com/rachel_os/fastsearch/searcher/model"
	"github.com/gin-gonic/gin"
)

// NegaTiveQuery 负面词查询
func QueryNegative(c *gin.Context) {
	q := c.Query("q")
	neg := &model.NegSearch{}
	if err := c.ShouldBindJSON(&neg); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	if neg.Query == "" {
		neg.Query = q
	}
	dbName := c.Query("database")
	data, err := srv.Negative.QueryNegative(dbName, neg)
	if err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	ResponseSuccessWithData(c, data)
}
func HasNegative(c *gin.Context) {
	q := c.Query("q")
	dbName := c.Query("database")
	neg := &model.NegSearch{}
	if neg.Query == "" {
		neg.Query = q
	}
	a, err := srv.Negative.HasNegative(dbName, neg)
	if err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	ResponseSuccessWithData(c, a)
}
func AllKeys(c *gin.Context) {
	q := c.Query("q")
	neg := &model.NegSearch{}
	c.ShouldBindJSON(&neg)
	if neg.Query == "" {
		neg.Query = q
	}
	dbName := c.Query("database")
	data, err := srv.Negative.AllKeys(dbName, neg)
	if err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	ResponseSuccessWithData(c, data)
}

// NegaTive_Add 导入负面词
func NegativeKeys(c *gin.Context) {
	neg := &model.NegativeDoc{}
	if err := c.ShouldBindJSON(&neg); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	dbName := c.Query("database")
	id, err := srv.Negative.AddNegative(dbName, neg)
	if err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	ResponseSuccessWithData(c, id)
}
func NegativeBatch(c *gin.Context) {
	negs := make([]model.NegativeDoc, 0)
	if err := c.BindJSON(&negs); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	dbName := c.Query("database")
	ids, err := srv.Negative.BatchAddNegative(dbName, &negs)
	if err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	ResponseSuccessWithData(c, ids)
}
func NegativeRemove(c *gin.Context) {
	neg := &model.RemoveNegativeModel{}
	if err := c.ShouldBindJSON(&neg); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	dbName := c.Query("database")
	if err := srv.Negative.RemoveNegative(dbName, neg); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}

	ResponseSuccessWithData(c, nil)
}
