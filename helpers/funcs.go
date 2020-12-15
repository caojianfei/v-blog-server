package helpers

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"v-blog/databases"
	"v-blog/models"
)

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 批量获取图片地址
func BatchGetImageUrlsByMd5(md5Arr []string) (map[string]string, error) {
	var files []models.File
	databases.DB.Where("md5 IN (?)", md5Arr).Find(&files)

	if len(files) == 0 {
		return nil, errors.New("文件不存在")
	}

	list := make(map[string]string)
	for _, f := range files {
		url, err := f.Url()
		if err != nil {
			list[f.Md5] = ""
		} else {
			list[f.Md5] = url
		}
	}

	return list, nil
}

func SingleGetImageUrlByMd5(md5 string) string {
	if md5 == "" {
		 return ""
	}
	var file models.File
	databases.DB.Where("md5 = ?", md5).First(&file)

	url, err := file.Url()
	if err != nil {
		return ""
	}

	return url
}

// 切片去重
func UniqueSlice(val interface{}) (interface{}, error) {

	s, ok := IsSlice(val)
	if !ok {
		return nil, errors.New("the param`s type is not slice")
	}

	us := make([]interface{}, 0, s.Len())
	for i := 0; i < s.Len(); i++ {
		repeat := false
		v := s.Index(i).Interface()
		for _, item := range us {
			if v == item {
				repeat = true
			}
		}
		if !repeat {
			us = append(us, v)
		}
	}

	return us, nil
}

// 判断参数是否是 slice
func IsSlice(s interface{}) (reflect.Value, bool) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Slice {
		return val, true
	}

	return val, false
}

func SliceColumn(s interface{}, field string) ([]interface{}, error) {
	vs, ok := IsSlice(s)
	if !ok {
		return nil, errors.New("the param`s type is not slice")
	}

	result := make([]interface{}, vs.Len())
	for i := 0; i < vs.Len(); i++ {
		val := vs.Index(i)
		// 非结构体
		if val.Kind() != reflect.Struct {
			return nil, errors.New("the item of slice is not struct")
		}
		fieldVal := val.FieldByName(field)
		if !fieldVal.IsValid() {
			return nil, errors.New(fmt.Sprintf("the item  of slice`s field is not valid. field name: [%s]", field))
		}

		result = append(result, fieldVal.Interface())
	}

	return result, nil
}

func SliceColumnUnique(s interface{}, field string) (interface{}, error) {
	columns, err := SliceColumn(s, field)
	if err != nil {
		return nil, err
	}

	uniqueColumns, err := UniqueSlice(columns)
	if err != nil {
		return nil, err
	}

	return uniqueColumns, nil
}
