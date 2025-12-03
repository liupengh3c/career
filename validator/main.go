package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// User contains user information
type User struct {
	FirstName      string     `validate:"required,min=2,max=10"`
	LastName       string     `validate:"required len=3"`
	Age            uint8      `json:"age" validate:"gte=0,lte=130"`
	Email          string     `json:"email" validate:"required,email"`
	Gender         string     `json:"gender" validate:"oneof=male female prefer_not_to"`
	FavouriteColor string     `json:"favourite_color" validate:"iscolor"`          // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `json:"addresses" validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `json:"street" validate:"required"`
	City   string `json:"city" validate:"required"`
	Planet string `json:"planet" validate:"required"`
	Phone  string `json:"phone" validate:"required"`
}

// CamelToSnake 转换驼峰为下划线命名
func CamelToSnake(s string) string {
	// 处理首字母大写的情况
	snake := strings.ToLower(s)

	// 正则1：在小写字母后的大写字母前插入下划线
	re1 := regexp.MustCompile(`([a-z])([A-Z])`)
	snake = re1.ReplaceAllString(snake, "$1_$2")
	// 正则2：处理数字或小写字母后的大写字母
	re2 := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake = re2.ReplaceAllString(snake, "$1_$2")

	return snake
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func main() {

	validate = validator.New(validator.WithRequiredStructEnabled())

	validateStruct()
	validateVariable()
}

func validateStruct() {

	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
	}

	user := &User{
		FirstName:      "BB",
		LastName:       "Smith",
		Age:            135,
		Gender:         "male",
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#000-",
		Addresses:      []*Address{address},
	}

	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(user)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			fmt.Println(err)
			return
		}

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			fmt.Println("验证失败，发现以下错误：")
			for _, e := range validateErrs {
				// 获取字段名（使用 StructField 获取结构体字段名）
				fieldName := e.StructField()
				if fieldName == "" {
					fieldName = e.Field()
				}
				fieldName = CamelToSnake(fieldName)
				// 获取错误原因
				var reason string
				switch e.Tag() {
				case "required":
					reason = "该字段为必填项"
				case "email":
					reason = "邮箱格式不正确"
				case "gte":
					reason = fmt.Sprintf("值必须大于等于 %s", e.Param())
				case "lte":
					reason = fmt.Sprintf("值必须小于等于 %s", e.Param())
				case "oneof":
					reason = fmt.Sprintf("值必须是以下之一: %s", e.Param())
				case "iscolor", "hexcolor", "rgb", "rgba", "hsl", "hsla":
					reason = "颜色格式不正确"
				default:
					reason = fmt.Sprintf("验证失败: %s", e.Tag())
				}

				fmt.Printf("  字段: %s, 错误: %s (当前值: %v)\n", fieldName, reason, e.Value())
			}
		}

		// from here you can create your own error messages in whatever language you wish
		return
	}

	// save user to database
}

func validateVariable() {

	myEmail := "joeybloggs.gmail.com"

	errs := validate.Var(myEmail, "required,email")

	if errs != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(errs, &validateErrs) {
			fmt.Println("变量验证失败：")
			for _, e := range validateErrs {
				var reason string
				switch e.Tag() {
				case "required":
					reason = "该字段为必填项"
				case "email":
					reason = "邮箱格式不正确"
				default:
					reason = fmt.Sprintf("验证失败: %s", e.Tag())
				}
				fmt.Printf("  错误: %s (当前值: %v)\n", reason, e.Value())
			}
		} else {
			fmt.Println("验证错误:", errs)
		}
		return
	}

	// email ok, move on
}
