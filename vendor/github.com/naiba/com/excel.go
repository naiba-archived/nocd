/*
 * Copyright (c) 2017, 奶爸<1@5.nu>
 * All rights reserved.
 */

package com

import (
	"github.com/xuri/excelize"
	"reflect"
	"strconv"
)

//创建带标题的表格
func NewTitledExcel(mStrut interface{}) *excelize.File {
	excel := excelize.NewFile()
	keys := reflect.TypeOf(mStrut)
	values := reflect.ValueOf(mStrut)
	for i := 0; i < keys.NumField(); i++ {
		var name = keys.Field(i).Name
		if values.Field(i).Kind() == reflect.Slice {
			name += "[]"
		}
		excel.SetCellValue("sheet1", GetColumnName(i)+"1", name)
	}
	return excel
}

//获取列名，从0开始获取
func GetColumnName(index int) (columnName string) {
	//A  B  C  AA  AB  AC  BA  BB  BC
	//0  1  2  3   4   5   6   7   8
	columnName = ""
	if index > 25 {
		if index%25 == 0 {
			columnName = string(rune(index/26+64)) + GetColumnName(26-index/25)
		} else {
			columnName = string(rune(index/26+64)) + GetColumnName(index%26)
		}
	} else {
		columnName = string(rune(index + 65))
	}
	return
}

//将结构体写入表格
func WriteStrutToExcel(xlsx *excelize.File, strucd interface{}, line int) {
	keys := reflect.TypeOf(strucd)
	values := reflect.ValueOf(strucd)
	for i := 0; i < keys.NumField(); i++ {
		var val string
		if values.Field(i).Kind() == reflect.Slice {
			for j := 0; j < values.Field(i).Len(); j++ {
				val += values.Field(i).Index(j).String() + "^"
			}
		} else {
			val = values.Field(i).String()
		}
		xlsx.SetCellValue("sheet1", GetColumnName(i)+strconv.Itoa(line), val)
	}
}
