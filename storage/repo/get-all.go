package repo

import "github.com/SaidakbarPardaboyev/get-all-from-ucode/storage/inner"

type GetAllI interface {
	Filter(filter map[string]interface{}) *inner.GetAll
	Sort(sort map[string]interface{}) *inner.GetAll
	Limit(limit int64) *inner.GetAll
	Skip(skip int64) *inner.GetAll
	Pipeline(pipeline []map[string]any) *inner.GetAll
	Count() (int64, error)
	Exec() ([]map[string]interface{}, error)
}
