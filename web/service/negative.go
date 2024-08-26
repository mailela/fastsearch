package service

import (
	"errors"

	"gitee.com/rachel_os/fastsearch/global"
	"gitee.com/rachel_os/fastsearch/negative"
	"gitee.com/rachel_os/fastsearch/searcher"
	"gitee.com/rachel_os/fastsearch/searcher/model"
)

type Negative struct {
	Container *searcher.Container
}

func NewNegative() *Negative {
	return &Negative{
		Container: global.Container,
	}
}

// AddNegative 添加索引
func (i *Negative) AddNegative(dbName string, request *model.NegativeDoc) (string, error) {
	id, err := i.Container.Neg.NegativeKeys(request)
	if err != nil {
		return id, errors.New("failed to generate negative key")
	}
	return id, nil
}
func (i *Negative) RemoveNegative(dbName string, request *model.RemoveNegativeModel) error {
	return i.Container.Neg.RemoveNegative(request.Id)
}
func (i *Negative) QueryNegative(dbName string, request *model.NegSearch) (*model.NegResult, error) {
	return i.Container.Neg.QueryNegative(request)
}
func (i *Negative) AllKeys(dbName string, request *model.NegSearch) ([]negative.KeyItem, error) {
	return i.Container.Neg.AllKeys(request)
}
func (i *Negative) HasNegative(dbName string, request *model.NegSearch) ([]negative.KeyItem, error) {
	a, _, c := i.Container.Neg.HasNegative(request.Query)
	return a, c
}

func (i *Negative) BatchAddNegative(dbName string, request *[]model.NegativeDoc) ([]string, error) {
	var ids []string
	ids, err := i.Container.Neg.BatchNegativeKeys(request)
	if err != nil {
		return ids, errors.New("failed to batch add negative")
	}
	return ids, nil
}
