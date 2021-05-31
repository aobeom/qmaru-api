package models

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// DBName 格式化表名
func DBName(s string) string {
	data := make([]string, 0)
	ok := false
	ru := []rune(s)

	for i := 0; i < len(ru); i++ {
		d := ru[i]
		if unicode.IsUpper(d) && ok {
			data = append(data, "_")
		}

		if string(d) != "_" {
			ok = true
		}
		data = append(data, string(d))
	}
	return strings.ToLower(strings.Join(data, ""))
}

// DBFiled 解析 model 字段
func DBFiled(reflectType reflect.Type, buffer *bytes.Buffer) {
	if reflectType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < reflectType.NumField(); i++ {
		jsonTag := reflectType.Field(i).Tag.Get("json")
		dbTag := reflectType.Field(i).Tag.Get("db")

		if jsonTag == "" && dbTag == "" {
			DBFiled(reflectType.Field(i).Type, buffer)
			continue
		}

		dbProfile := strings.Split(dbTag, ";")
		dbFiled := fmt.Sprintf("%s %s", jsonTag, strings.Join(dbProfile, " "))
		buffer.WriteString(dbFiled)
		buffer.WriteString(",")
	}
}
