package commom

import (
	"gorm.io/gorm"
	"reflect"
)

type SelectOption struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

type TreeSelectOption struct {
	Label    string `json:"label"`
	Value    int    `json:"value"`
	ParentID int    `json:"parent_id"`
}

type CommonResult struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func CommonResultSuccess(data interface{}) CommonResult {
	return CommonResult{
		Code: 0,
		Data: data,
		Msg:  "",
	}
}

func CommonResultFailed(e error) CommonResult {
	return CommonResult{
		Code: -1,
		Data: nil,
		Msg:  e.Error(),
	}
}

func CommonResultAuto(data interface{}, e error) CommonResult {
	if e != nil {
		return CommonResultFailed(e)
	} else {
		return CommonResultSuccess(data)
	}
}

type Pageable interface {
	Pageable() (int, int)
}

func PaginateQuery(paginate Pageable) func(db *gorm.DB) *gorm.DB {
	limit, offset := paginate.Pageable()
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit).Offset(offset)
	}
}

func Map2ID(group interface{}) []int {
	ids := make([]int, 0)
	val := reflect.ValueOf(group)
	if val.Kind() != reflect.Slice {
		return ids
	}
	for i := 0; i < val.Len(); i++ {
		ids = append(ids, val.Index(i).FieldByName("ID").Interface().(int))
	}
	return ids
}
