package repo

import (
	"github.com/SaidakbarPardaboyev/get-all-from-ucode/pkg"
	"github.com/SaidakbarPardaboyev/get-all-from-ucode/storage/inner"
)

type ItemsI interface {
	GetAll() GetAllI
	MultipleUpdate() MultipleUpdateI
}

type APIItem struct {
	Collection string
	Config     *pkg.InnerConfig
}

func (i *APIItem) GetAll() GetAllI {
	return inner.NewGetAllRepo(i.Collection, i.Config)
}

func (i *APIItem) MultipleUpdate() MultipleUpdateI {
	return inner.NewMultipleUpdate(i.Collection, i.Config)
}
