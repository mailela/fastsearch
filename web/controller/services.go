package controller

import (
	service2 "gitee.com/rachel_os/fastsearch/web/service"
)

var srv *Services

type Services struct {
	Base     *service2.Base
	Index    *service2.Index
	Negative *service2.Negative
	Database *service2.Database
	Word     *service2.Word
}

func NewServices() {
	srv = &Services{
		Base:     service2.NewBase(),
		Index:    service2.NewIndex(),
		Negative: service2.NewNegative(),
		Database: service2.NewDatabase(),
		Word:     service2.NewWord(),
	}
}
