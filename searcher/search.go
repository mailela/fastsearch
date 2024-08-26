package searcher

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"

	"gitee.com/rachel_os/fastsearch/negative"
	"gitee.com/rachel_os/fastsearch/searcher/exp"
	"gitee.com/rachel_os/fastsearch/searcher/model"
	"gitee.com/rachel_os/fastsearch/searcher/pagination"
	"gitee.com/rachel_os/fastsearch/searcher/sorts"
	"gitee.com/rachel_os/fastsearch/searcher/utils"
)

func (e *Engine) MultiSearch(request *model.SearchRequest) (any, error) {
	//等待搜索初始化完成
	e.Wait()
	if request.Negative.Query {
		neg_words, no_pass, _ := negative.Neg.HasNegative(request.Query)
		if no_pass {
			return neg_words, errors.New("包含负面词")
		}
	}

	//分词搜索
	words := e.Tokenizer.Cut(request.Query)

	fastSort := &sorts.FastSort{
		IsDebug: e.IsDebug,
		Order:   request.Order,
	}

	_time := utils.ExecTime(func() {

		base := len(words)
		wg := &sync.WaitGroup{}
		wg.Add(base)

		for _, word := range words {
			go e.processKeySearch(word, fastSort, wg)
		}
		wg.Wait()
	})
	if e.IsDebug {
		log.Println("搜索时间:", _time, "ms")
	}
	// 处理分页
	request = request.GetAndSetDefault()

	//计算交集得分和去重
	fastSort.Process()

	wordMap := make(map[string]bool)
	for _, word := range words {
		wordMap[word] = true
	}

	//读取文档
	var result = &model.SearchResult{
		Page:  request.Page,
		Limit: request.Limit,
		Words: words,
	}

	t, err := utils.ExecTimeWithError(func() error {
		pager := new(pagination.Pagination)
		result.Total = fastSort.Count()
		//设置总页数
		result.PageCount = pager.PageCount
		pager.Init(request.Limit, result.Total)

		//读取单页的id
		if pager.PageCount != 0 {
			// 生成计算表达式
			filter_Exp, _ := exp.NewEvaluableExpression(request.FilterExp)
			score_Exp, _ := exp.NewEvaluableExpression(request.ScoreExp)
			start, end := pager.GetPage(request.Page)
			var resultItems = make([]model.SliceItem, 0)
			fastSort.GetAll(&resultItems, start, end)
			count := len(resultItems)

			//获取文档
			result.Documents = make([]model.ResponseDoc, count)

			var (
				tmp_docs   []model.ResponseDoc
				found      uint64
				pagesize   uint64
				page_count int
			)
			pagesize = uint64(request.Limit)
			page_count = result.Total
			result.PageCount = pager.PageCount
			if (page_count - start) < int(pagesize) {
				result.PageCount = int(math.Ceil(float64(end / int(pagesize))))
			}
			for index := start; index < page_count; index++ {
				if found >= pagesize {
					result.PageCount = result.Page + 1
					break
				}

				item := fastSort.Get(index)
				doc, _ := e.Get(item, request, &wordMap)
				parameters := utils.Obj2Map(doc)

				// 计算分数
				if request.ScoreExp != "" {
					val, err := score_Exp.Evaluate(parameters)
					if err != nil {
						log.Printf("表达式执行'%v'错误: %v 值内容: %v", request.ScoreExp, err, parameters)
					} else if val != nil {
						doc.Score = int(val.(float64))
					}
				}
				// 过滤结果
				if request.FilterExp != "" {
					val, err := filter_Exp.Evaluate(parameters)
					if err != nil {
						log.Printf("表达式执行'%v'错误: %v 值内容: %v", request.ScoreExp, err, parameters)
					} else if val != nil {
						if val.(bool) {
							tmp_docs = append(tmp_docs, doc)
						} else {
							result.Total--
							continue
						}
					}
				} else {
					tmp_docs = append(tmp_docs, doc)
				}
				if request.Negative.Content {
					text := fmt.Sprintf("%s", doc.Text)
					_, no_pass, _ := negative.Neg.HasNegative(text)
					if no_pass {
						s := len(tmp_docs) - 1
						tmp_docs = append(tmp_docs[:s], tmp_docs[s+1:]...)
						result.Total--
						continue
					}
				}
				found++
			}

			result.Documents = tmp_docs
			pager.Total = result.Total
			if request.Order == "desc" {
				sort.Sort(sort.Reverse(model.ResponseDocSort(result.Documents)))
			} else {
				sort.Sort(model.ResponseDocSort(result.Documents))
			}

			if request.ScoreExp != "" || request.FilterExp != "" {
				// 取出page
				start, found := pager.GetPage(request.Page)
				result.Documents = result.Documents[start:found]
			}
		}
		return nil
	})

	if e.IsDebug {
		log.Println("处理数据耗时：", _time, "ms")
	}
	if err != nil {
		return nil, err
	}
	result.Time = _time + t

	return result, nil
}
