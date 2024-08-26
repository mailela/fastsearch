package negative

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gitee.com/rachel_os/fastsearch/searcher/model"
	"gitee.com/rachel_os/fastsearch/searcher/pagination"
	"gitee.com/rachel_os/fastsearch/searcher/storage"
	"gitee.com/rachel_os/fastsearch/searcher/utils"
	"gitee.com/rachel_os/fastsearch/searcher/words"
)

var (
	Keys []KeyItem
	Neg  *Engine
)

type KeyItem struct {
	Key  string
	Type string
}
type Engine struct {
	NegativeStorage *storage.LeveldbStorage //存储负向文档的storage
	Tokenizer       *words.Tokenizer
	NegativeCount   int //负向文档数量
}

func NewNegative(path string, tokenizer *words.Tokenizer) *Engine {
	neg, err := storage.NewStorage(path, 10)
	if err != nil {
		panic(err)
	}
	Neg = &Engine{
		NegativeStorage: neg,
		Tokenizer:       tokenizer,
		NegativeCount:   int(neg.GetCount()),
	}
	return Neg
}

func (e *Engine) QueryNegative(s *model.NegSearch) (*model.NegResult, error) {
	var (
		items storage.Items
		count int64
		doc   = model.NegativeDoc{}
	)
	result := &model.NegResult{Time: 0}
	s = s.GetAndSetDefault()
	docStorage := e.NegativeStorage
	start := (s.Page - 1) * s.Limit
	_time := utils.ExecTime(func() {
		count, items = docStorage.Search(func(item storage.Item) any {
			utils.Decoder(item.Value, &doc)
			return doc
		}, &storage.SearchOption{
			Start:     int64(start),
			End:       int64(start + s.Limit),
			FilterExp: s.FilterExp,
			Order:     s.Order,
		})
	})
	result.Time = _time
	pager := new(pagination.Pagination)
	pager.Init(s.Limit, int(count))

	result.Total = pager.Total
	result.PageCount = pager.PageCount

	docs := make([]model.NegativeDoc, 0)
	for _, item := range items {
		utils.Decoder(item.Value, &doc)
		// docs = append(docs, model.NegativeDoc{Keys: doc.Keys, Title: doc.Title, Id: doc.Id, Time: doc.Time})
		docs = append(docs, doc)
	}
	result.Documents = docs
	return result, nil
}

// AllKeys 函数从负向搜索的文档中获取所有关键词并返回
// Engine是调用该方法的类型
// docs是指向model.NegSearch的指针，表示负向搜索的文档
// 返回值为一个KeyItem类型的切片和error类型
func (e *Engine) AllKeys(docs *model.NegSearch) ([]KeyItem, error) {
	var (
		result []KeyItem
		count  int64
		items  storage.Items
		doc    = model.NegativeDoc{}
	)
	result = make([]KeyItem, 0)
	docStorage := e.NegativeStorage
	count, items = docStorage.GetAll(0, 0, "asc", nil)
	if count <= 0 {
		return result, nil
	}
	for _, item := range items {
		utils.Decoder(item.Value, &doc)
		for _, key := range doc.Keys {
			key = strings.TrimSpace(key)
			if key == "" {
				continue
			}
			result = append(result, KeyItem{Key: key, Type: doc.Type})
		}
	}
	Keys = result
	return result, nil
}

func (e *Engine) RemoveNegative(id string) error {
	//移除
	docStorage := e.NegativeStorage
	key := []byte(id)

	//数量增加
	if !docStorage.Has(key) {
		return errors.New(fmt.Sprintf("negative %s not exist", id))
	}
	err := docStorage.Delete(key)
	if err != nil {
		e.NegativeCount--
	}
	return nil
}
func (e *Engine) NegativeKeys(doc *model.NegativeDoc) (string, error) {
	doc.Time = time.Now().String()
	docStorage := e.NegativeStorage

	doc.Keys = e.Tokenizer.Split(doc.Text)
	doc.Text = strings.Join(doc.Keys, "\n")
	if doc.Id == "" {
		doc.Id = utils.Sha1Hash(doc.Text)
	}
	if doc.Type == "" {
		doc.Type = "0"
	}
	key := []byte(doc.Id)
	doc.Time = utils.Now()
	//数量增加
	if !docStorage.Has(key) {
		e.NegativeCount++
	}
	if doc.Title == "" {
		doc.Title = doc.Text
	}
	//设置到id和key的映射中
	docStorage.Set(key, utils.Encoder(doc))
	return doc.Id, nil
}

// BatchNegativeKeys 函数批量处理NegativeDoc类型指针切片，将其存储在NegativeStorage中，并返回存储成功的id列表和错误信息
func (e *Engine) BatchNegativeKeys(docs *[]model.NegativeDoc) ([]string, error) {
	docStorage := e.NegativeStorage
	var (
		items storage.Items
	)
	for _, doc := range *docs {

		doc.Keys = e.Tokenizer.Split(doc.Text)
		doc.Text = strings.Join(doc.Keys, "\n")

		if doc.Id == "" {
			doc.Id = utils.Sha1Hash(doc.Text)
		}
		if doc.Type == "" {
			doc.Type = "0"
		}
		key := []byte(doc.Id)
		if doc.Title == "" {
			doc.Title = doc.Text
		}
		doc.Time = utils.Now()
		items = append(items, storage.Item{Key: key, Value: utils.Encoder(doc)})
	}
	ids, err := docStorage.BatchSet(&items)
	return ids, err
}

// 检测是否在关键字中
func (e *Engine) HasNegative(text string) ([]KeyItem, bool, error) {
	var (
		result []KeyItem
	)
	for _, v := range Keys {
		exp_str := v.Key
		rel := utils.FindByExp(text, exp_str)
		if rel {
			result = append(result, v)
			return result, true, nil
		}
	}
	return result, false, nil
}
