package slice

import (
	"errors"
	"fmt"
	"reflect"
)

type Element interface{}

type Slice []Element

// 将任意类型切片转换成 Slice
func ToSlice(s interface{}) Slice {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Slice {
		panic("the param is not slice")
	}

	result := make(Slice, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		result = append(result, Element(val.Index(i).Interface()))
	}

	return result
}

// 将结构体类型切片指定字段取出，组成新的 Slice
func (s Slice) Column(field string) Slice {
	result := make(Slice, 0, len(s))
	for _, val := range s{
		reflectVal := reflect.ValueOf(val)
		if reflectVal.Kind() != reflect.Struct {
			panic("the slice item is not struct")
		}

		fieldValue := reflectVal.FieldByName(field)
		if !fieldValue.IsValid() {
			panic(fmt.Sprintf("field [%s] is not exist", field))
		}

		result = append(result, Element(fieldValue.Interface()))
	}

	return result
}

// 切片去重
func (s Slice) Unique() Slice {
	sl := len(s)
	if sl == 0 {
		return s
	}
	sMap := make(map[Element]Element)
	for i := 0; i < sl; i++ {
		item := s[i]
		itemType := reflect.TypeOf(item).Kind()
		// slice, function, map 不能使用 == 比较，所以无法去重
		if itemType == reflect.Slice || itemType == reflect.Func || itemType == reflect.Map {
			panic("[slice, function, map] slice can not use this function")
		}
		if _, ok := sMap[s[i]]; !ok {
			sMap[s[i]] = s[i]
		}
	}

	result := make(Slice, 0, sl)
	for _, val := range sMap {
		result = append(result, val)
	}

	return result
}

// 转换成 int 型切片
func (s Slice) CoverToInt() ([]int, error) {
	result := make([]int, 0, len(s))
	for key, val := range s {
		if v, ok := val.(int); ok {
			result = append(result, v)
		} else {
			return []int{}, errors.New(fmt.Sprintf("the value [%v] is not int. key: [%d].", val, key))
		}
	}

	return result, nil
}

// 转换成 string 型切片
func (s Slice) CovertToString() ([]string, error) {
	result := make([]string, 0, len(s))
	for key, val := range s {
		if v, ok := val.(string); ok {
			result = append(result, v)
		} else {
			return []string{}, errors.New(fmt.Sprintf("the value [%v] is not string. key: [%d].", val, key))
		}
	}

	return result, nil
}
