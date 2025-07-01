package testhelper

import (
	"strconv"

	"github.com/buger/jsonparser"

	string_helper "financing-offer/pkg/string-helper"
)

func GetArrayLength(data []byte, key ...string) int {
	val, typ, _, err := jsonparser.Get(data, key...)
	if err != nil || typ != jsonparser.Array {
		return 0
	}
	length := 0
	_, _ = jsonparser.ArrayEach(
		val, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			length++
		},
	)
	return length
}

func GetArrayInt(data []byte, key ...string) []int64 {
	res := make([]int64, 0)
	_, _ = jsonparser.ArrayEach(
		data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			val, _ := strconv.ParseInt(string(value), 10, 64)
			res = append(res, val)
		}, key...,
	)
	return res
}

func GetArrayString(data []byte, key ...string) []string {
	res := make([]string, 0)
	_, _ = jsonparser.ArrayEach(
		data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			res = append(res, string_helper.BytesToString(value))
		}, key...,
	)
	return res
}

func GetString(data []byte, key ...string) string {
	res, _ := jsonparser.GetString(data, key...)
	return res
}

func GetInt(data []byte, key ...string) int64 {
	res, _ := jsonparser.GetInt(data, key...)
	return res
}

func GetFloat(data []byte, key ...string) float64 {
	res, _ := jsonparser.GetFloat(data, key...)
	return res
}

func GetBoolean(data []byte, key ...string) bool {
	res, _ := jsonparser.GetBoolean(data, key...)
	return res
}
