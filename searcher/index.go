package searcher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"gitee.com/rachel_os/fastsearch/negative"
	"gitee.com/rachel_os/fastsearch/searcher/arrays"
	"gitee.com/rachel_os/fastsearch/searcher/exp"
	"gitee.com/rachel_os/fastsearch/searcher/model"
	"gitee.com/rachel_os/fastsearch/searcher/pagination"
	"gitee.com/rachel_os/fastsearch/searcher/sorts"
	"gitee.com/rachel_os/fastsearch/searcher/storage"
	"gitee.com/rachel_os/fastsearch/searcher/utils"
)

func (e *Engine) IndexDocument(doc *model.IndexDoc) error {
	doc.Time = utils.Now()
	//数量增加
	e.addDocumentWorkerChan[e.getShard(doc.Id)] <- doc
	return nil
}

// GetQueue 获取队列剩余
func (e *Engine) GetQueue() int {
	total := 0
	for _, v := range e.addDocumentWorkerChan {
		total += len(v)
	}
	return total
}

// DocumentWorkerExec 添加文档队列
func (e *Engine) DocumentWorkerExec(worker chan *model.IndexDoc) {
	for {
		doc := <-worker
		e.AddDocument(doc)
	}
}

// getShardKey 计算索引分布在哪个文件块
func (e *Engine) getShard(id string) int {
	// return int(id % uint32(e.Shard))
	sid := e.getShardKeyByWord(id)
	// fmt.Println("sid:", sid)
	return sid
}

func (e *Engine) getShardKeyByWord(word string) int {
	return int(utils.StringToInt(word) % uint32(e.Shard))
}

func (e *Engine) InitOption(option *Option) {

	if option == nil {
		//默认值
		option = e.GetOptions()
	}
	e.Option = option
	//shard默认值
	if e.Shard <= 0 {
		e.Shard = 10
	}
	if e.BufferNum <= 0 {
		e.BufferNum = 1000
	}
	//初始化其他的
	e.Init()

}

func (e *Engine) getFilePath(fileName string) string {
	return e.IndexPath + string(os.PathSeparator) + fileName
}

func (e *Engine) GetOptions() *Option {
	return &Option{
		DocIndexName:      "docs",
		InvertedIndexName: "inverted_index",
		PositiveIndexName: "positive_index",
	}
}

// AddDocument 分词索引
func (e *Engine) AddDocument(index *model.IndexDoc) {
	//等待初始化完成
	e.Wait()
	text := index.Text

	splitWords := e.Tokenizer.Cut(text)

	id := index.Id
	// 检查是否需要更新倒排索引 words变更/id不存在
	inserts, needUpdateInverted := e.optimizeIndex(id, splitWords)
	// 将新增的word剔出单独处理，减少I/O操作
	if needUpdateInverted {
		for _, word := range inserts {
			e.addInvertedIndex(word, id)
		}
	}

	// TODO: 是否需要更新正排索引 - 检测document变更
	e.addPositiveIndex(index, splitWords)
}

// 添加倒排索引
func (e *Engine) addInvertedIndex(word string, id string) {
	e.Lock()
	defer e.Unlock()

	shard := e.getShardKeyByWord(word)

	s := e.invertedIndexStorages[shard]

	//string作为key
	key := []byte(word)

	//存在
	//添加到列表
	buf, find := s.Get(key)
	// ids := make([]uint32, 0)
	ids := make([]string, 0)
	if find {
		utils.Decoder(buf, &ids)
	}

	if !arrays.ArrayStringExists(ids, id) {
		ids = append(ids, id)
	}

	s.Set(key, utils.Encoder(ids))
}

// 移除删去的词
func (e *Engine) optimizeIndex(id string, newWords []string) ([]string, bool) {
	// 判断id是否存在
	e.Lock()
	defer e.Unlock()

	// 计算差值
	removes, inserts, changed := e.getDifference(id, newWords)
	if changed {
		if removes != nil && len(removes) > 0 {
			// 移除正排索引
			for _, word := range removes {
				e.removeIdInWordIndex(id, word)
			}
		}
	}
	return inserts, changed
}

func (e *Engine) removeIdInWordIndex(id string, word string) {

	shard := e.getShardKeyByWord(word)

	wordStorage := e.invertedIndexStorages[shard]

	//string作为key
	key := []byte(word)

	buf, found := wordStorage.Get(key)
	if found {
		// ids := make([]uint32, 0)
		ids := make([]string, 0)
		utils.Decoder(buf, &ids)

		//移除
		index := arrays.FindString(ids, id)
		if index != -1 {
			ids = utils.DeleteStringArray(ids, index)
			if len(ids) == 0 {
				err := wordStorage.Delete(key)
				if err != nil {
					panic(err)
				}
			} else {
				wordStorage.Set(key, utils.Encoder(ids))
			}
		}
	}

}

// 计算差值
// @return []string: 需要删除的词
// @return bool    : words出现变更返回true，否则返回false
func (e *Engine) getDifference(id string, newWords []string) ([]string, []string, bool) {
	shard := e.getShard(id)
	wordStorage := e.positiveIndexStorages[shard]
	// key := utils.Uint32ToBytes(id)
	key := []byte(id)
	buf, found := wordStorage.Get(key)
	if found {
		oldWords := make([]string, 0)
		utils.Decoder(buf, &oldWords)

		// 计算需要移除的
		removes := make([]string, 0)
		for _, word := range oldWords {
			// 旧的在新的里面不存在，就是需要移除的
			if !arrays.ArrayStringExists(newWords, word) {
				removes = append(removes, word)
			}
		}
		// 计算需要新增的
		inserts := make([]string, 0)
		for _, word := range newWords {
			if !arrays.ArrayStringExists(oldWords, word) {
				inserts = append(inserts, word)
			}
		}
		if len(removes) != 0 || len(inserts) != 0 {
			return removes, inserts, true
		}
		// 没有改变
		return removes, inserts, false
	}
	// id不存在，相当于insert
	return nil, newWords, true
}

// 添加正排索引 id=>keys id=>doc
func (e *Engine) addPositiveIndex(index *model.IndexDoc, keys []string) {
	e.Lock()
	defer e.Unlock()

	// key := utils.Uint32ToBytes(index.Id)
	key := []byte(index.Id)
	shard := e.getShard(index.Id)
	docStorage := e.docStorages[shard]

	//id和key的映射
	positiveIndexStorage := e.positiveIndexStorages[shard]

	doc := &model.StorageIndexDoc{
		IndexDoc: index,
		Keys:     keys,
	}
	//存储id和key以及文档的映射
	if !docStorage.Has(key) {
		e.DocumentCount++
	}
	docStorage.Set(key, utils.Encoder(doc))
	//设置到id和key的映射中
	positiveIndexStorage.Set(key, utils.Encoder(keys))
}

// MultiSearch 多线程搜索

func (e *Engine) MultiSearch1(request *model.SearchRequest) (any, error) {
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
		max_count := fastSort.Count()
		pager := new(pagination.Pagination)
		if request.FilterExp != "" || request.ScoreExp != "" {
			if request.MaxLimit == 0 {
				max_count = 1000
			} else if request.MaxLimit == -1 {
				max_count = fastSort.Count()
			}
			if fastSort.Count() < max_count {
				max_count = fastSort.Count()
			}
		}
		result.Total = max_count
		pager.Init(request.Limit, max_count)
		//设置总页数
		result.PageCount = pager.PageCount

		//读取单页的id
		if pager.PageCount != 0 {

			start, end := pager.GetPage(request.Page)
			if request.ScoreExp != "" || request.FilterExp != "" {
				// 分数表达式或过滤表达式不为空,获取所有的数据
				start, end = 0, pager.Total
			}
			var resultItems = make([]model.SliceItem, 0)
			fastSort.GetAll(&resultItems, start, end)
			count := len(resultItems)

			//获取文档
			result.Documents = make([]model.ResponseDoc, count)

			// 生成计算表达式
			filter_Exp, _ := exp.NewEvaluableExpression(request.FilterExp)
			score_Exp, _ := exp.NewEvaluableExpression(request.ScoreExp)
			tmp_docs := make([]model.ResponseDoc, 0)

			var (
				found    uint64 = 0
				pagesize uint64
			)
			pagesize = uint64(request.Limit)
			wg := new(sync.WaitGroup)
			wg.Add(count)
			for index, item := range resultItems {

				if found > pagesize {
					break
				}

				e.getDocument(item, &result.Documents[index], request, &wordMap, wg)
				// }
				// wg.Wait()

				// // 根据表达式进行过滤
				// for i, doc := range result.Documents {
				doc := result.Documents[index]
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
				start, end := pager.GetPage(request.Page)
				result.Documents = result.Documents[start:end]
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
func (e *Engine) Get(item model.SliceItem, request *model.SearchRequest, wordMap *map[string]bool) (model.ResponseDoc, error) {
	var (
		doc model.ResponseDoc
	)
	buf := e.GetDocById(string(item.Id))
	doc.Score = item.Score
	if buf != nil {
		//gob解析
		storageDoc := new(model.StorageIndexDoc)
		utils.Decoder(buf, &storageDoc)
		doc.Document = storageDoc.Document
		doc.Keys = storageDoc.Keys
		text := utils.RemoveHTMLTags(storageDoc.Text)
		title := utils.RemoveHTMLTags(storageDoc.Title)
		//处理关键词高亮
		highlight := request.Highlight
		if highlight != nil {
			//全部小写
			// text = strings.ToLower(text)
			//还可以优化，只替换击中的词
			for _, key := range storageDoc.Keys {
				if ok := (*wordMap)[key]; ok {
					key = utils.RemoveHTMLTags(key)
					text = strings.ReplaceAll(text, key, fmt.Sprintf("%s%s%s", highlight.PreTag, key, highlight.PostTag))
					title = strings.ReplaceAll(title, key, fmt.Sprintf("%s%s%s", highlight.PreTag, key, highlight.PostTag))
				}
			}
			//放置原始文本
			doc.OriginalText = storageDoc.Text
			doc.OriginalTitle = storageDoc.Title
		}
		doc.Text = text
		doc.Tags = storageDoc.Tags
		doc.Time = storageDoc.Time
		doc.Flag = storageDoc.Flag
		doc.Title = title
		doc.Id = string(item.Id)
	}
	return doc, nil
}
func (e *Engine) getDocument(item model.SliceItem, doc *model.ResponseDoc, request *model.SearchRequest, wordMap *map[string]bool, wg *sync.WaitGroup) {
	defer wg.Done()
	*doc, _ = e.Get(item, request, wordMap)
}

// processKeySearch 实现了Engine结构体中处理关键字搜索的功能
//
// 参数：
// word string - 待搜索的关键字
// fastSort *sorts.FastSort - 用于排序搜索结果的快速排序对象
// wg *sync.WaitGroup - 等待组对象，用于等待该协程执行完毕
//
// 返回值：
// 无返回值
func (e *Engine) processKeySearch(word string, fastSort *sorts.FastSort, wg *sync.WaitGroup) {
	defer wg.Done()

	shard := e.getShardKeyByWord(word)
	//读取id
	invertedIndexStorage := e.invertedIndexStorages[shard]
	key := []byte(word)

	buf, find := invertedIndexStorage.Get(key)
	if find {
		// ids := make([]uint32, 0)
		ids := make([]string, 0)
		//解码
		utils.Decoder(buf, &ids)
		fastSort.Add(&ids)
	}

}

// GetIndexCount 获取索引数量
func (e *Engine) GetIndexCount() int64 {
	var size int64
	for i := 0; i < e.Shard; i++ {
		size += e.invertedIndexStorages[i].GetCount()
	}
	return size
}

// GetDocumentCount 获取文档数量
func (e *Engine) GetDocumentCount() int64 {
	if e.DocumentCount == -1 {
		var count int64
		//使用多线程加速统计
		wg := sync.WaitGroup{}
		wg.Add(e.Shard)
		//这里的统计可能会出现数据错误，因为没加锁
		for i := 0; i < e.Shard; i++ {
			go func(i int) {
				count += e.docStorages[i].GetCount()
				wg.Done()
			}(i)
		}
		wg.Wait()
		e.DocumentCount = count
		return count + 1
	}

	return e.DocumentCount + 1
}

// GetDocById 通过id获取文档
func (e *Engine) GetDocById(id string) []byte {
	shard := e.getShard(id)
	// key := utils.Uint32ToBytes(id)
	key := []byte(id)
	buf, found := e.docStorages[shard].Get(key)
	if found {
		return buf
	}

	return nil
}

// RemoveIndex 根据ID移除索引
func (e *Engine) RemoveIndex(id string) error {

	//移除
	e.Lock()
	defer e.Unlock()

	shard := e.getShard(id)
	// key := utils.Uint32ToBytes(id)
	key := []byte(id)
	//关键字和Id映射
	//invertedIndexStorages []*storage.LeveldbStorage
	//ID和key映射，用于计算相关度，一个id 对应多个key
	ik := e.positiveIndexStorages[shard]
	keysValue, found := ik.Get(key)
	if !found {
		return errors.New(fmt.Sprintf("没有找到id=%s", id))
	}

	keys := make([]string, 0)
	utils.Decoder(keysValue, &keys)

	//符合条件的key，要移除id
	for _, word := range keys {
		e.removeIdInWordIndex(id, word)
	}

	//删除id映射
	err := ik.Delete(key)
	if err != nil {
		return errors.New(err.Error())
	}

	//文档仓
	err = e.docStorages[shard].Delete(key)
	if err != nil {
		return err
	}
	//减少数量
	e.DocumentCount--

	return nil
}

func (e *Engine) Close() {
	e.Lock()
	defer e.Unlock()

	for i := 0; i < e.Shard; i++ {
		e.invertedIndexStorages[i].Close()
		e.positiveIndexStorages[i].Close()
	}
}

// Drop 删除
func (e *Engine) Drop() error {
	if !e.AllowDrop {
		return errors.New("不允许删除")
	}
	e.Lock()
	defer e.Unlock()
	//删除文件
	if err := os.RemoveAll(e.IndexPath); err != nil {
		return err
	}

	//清空内存
	for i := 0; i < e.Shard; i++ {
		e.docStorages = make([]*storage.LeveldbStorage, 0)
		e.invertedIndexStorages = make([]*storage.LeveldbStorage, 0)
		e.positiveIndexStorages = make([]*storage.LeveldbStorage, 0)
	}

	return nil
}
