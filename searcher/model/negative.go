package model

import (
	"fmt"

	"gitee.com/rachel_os/fastsearch/searcher/utils"
)

type NegativeDoc struct {
	Id    string   `json:"id,omitempty"`
	Title string   `json:"title,omitempty"`
	Text  string   `json:"text,omitempty"`
	Type  string   `json:"type,omitempty"`
	Keys  []string `json:"keys,omitempty"`
	Time  string   `json:"time,omitempty"`
}
type RemoveNegativeModel struct {
	Id string `json:"id,omitempty"`
}

// NegRequest 搜索请求
type NegSearch struct {
	Query     string `json:"query,omitempty"`     // 搜索关键词
	Order     string `json:"order,omitempty"`     // 排序类型
	MaxLimit  int    `json:"maxlimit,omitempty"`  // 最大结果集限制，默认1000
	Flag      string `json:"flag,omitempty"`      // 搜索标记
	FilterExp string `json:"filterExp,omitempty"` // 过滤计算表达式
	ScoreExp  string `json:"scoreExp,omitempty"`  // 过滤计算表达式
	Page      int    `json:"page,omitempty"`      // 页码
	Limit     int    `json:"limit,omitempty"`     // 每页大小，最大1000，超过报错
}

func (s *NegSearch) GetAndSetDefault() *NegSearch {

	if s.Limit == 0 {
		s.Limit = 10
	}
	if s.Page == 0 {
		s.Page = 1
	}

	if s.Order == "" {
		s.Order = "desc"
	}
	if s.FilterExp == "" {
		if s.Query != "" {
			s.FilterExp = fmt.Sprintf("[text,title,time] Like '%s'", utils.RemoveSpace(s.Query))
		}
	}

	return s
}

// SearchResult 搜索响应
type NegResult struct {
	Time      float64       `json:"time,omitempty"`      //查询用时
	Total     int           `json:"total"`               //总数
	PageCount int           `json:"pageCount"`           //总页数
	Page      int           `json:"page,omitempty"`      //页码
	Limit     int           `json:"limit,omitempty"`     //页大小
	Documents []NegativeDoc `json:"documents,omitempty"` //文档
	Words     []string      `json:"words,omitempty"`     //搜索关键词
}
