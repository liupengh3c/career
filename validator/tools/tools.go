package tools

import (
	"fmt"
	"reflect"
	"strings"
)

type User struct {
	ID       int     `json:"id"`
	UserName string  `json:"user_name"`
	Email    string  `json:"email,omitempty"`
	Password string  `json:"-"` // 忽略字段
	Age      int     // 没有 json tag
	Addr     Address `json:"address"`
}

type Address struct {
	Addr  string `json:"addr"`
	Phone string `json:"phone"`
}

// 递归获取所有字段的 JSON tag
func GetJsonTagsRecursive(s interface{}) map[string]string {
	result := make(map[string]string)
	getJsonTagsRecursiveHelper(s, "", result)
	return result
}

func getJsonTagsRecursiveHelper(s interface{}, prefix string, result map[string]string) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue // 跳过被忽略的字段
		}

		fieldName := field.Name
		fullName := fieldName
		if prefix != "" {
			fullName = prefix + "." + fieldName
		}

		// 如果是嵌套结构体且不是匿名字段
		if fieldValue.Kind() == reflect.Struct && !field.Anonymous {
			if jsonTag != "" {
				// 使用 JSON tag 作为前缀
				jsonName := strings.Split(jsonTag, ",")[0]
				if jsonName != "" {
					getJsonTagsRecursiveHelper(fieldValue.Interface(), jsonName, result)
				} else {
					getJsonTagsRecursiveHelper(fieldValue.Interface(), fullName, result)
				}
			} else {
				getJsonTagsRecursiveHelper(fieldValue.Interface(), fullName, result)
			}
		} else {
			result[fullName] = strings.Split(jsonTag, ",")[0]
		}
	}
}

func main() {
	user := User{
		ID:       1,
		UserName: "john_doe",
		Email:    "john@example.com",
		Password: "secret",
		Age:      30,
	}

	tags := GetJsonTagsRecursive(user)
	for fieldName, jsonTag := range tags {
		fmt.Printf("字段: %-10s JSON Tag: %s\n", fieldName, jsonTag)
	}
}
