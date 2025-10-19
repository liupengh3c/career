package main

import (
	"fmt"
	"os"
	"text/template"
	"time"
)

type User struct {
	Name string
	Age  int
}

// 字符串反转函数
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// 日期格式化函数
func formatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
func main() {
	// TemplateVariables()
	FuncMapTemp()
}

func FuncMapTemp() {
	funcMap := template.FuncMap{
		"reverse": reverse,                             // 注册字符串反转函数
		"format":  formatDate,                          // 注册日期格式化函数
		"add":     func(a, b int) int { return a + b }, // 匿名函数
	}

	tmpl := template.Must(template.New("funcMap").Funcs(funcMap).Parse(`
        <h1>Hello {{.Name}}!</h1>
        <p>Reversed: {{ reverse .Name }}</p>
        <p>Date: {{ format .Now }}</p>
        <p>Sum: {{ add 10 20 }}</p>
    `))

	data := struct {
		Name string
		Now  time.Time
	}{"Go", time.Now()}

	tmpl.Execute(os.Stdout, data)
}

func RangeData() {
	items := []string{"apple", "banana", "orange"}
	tmpl, err := template.New("test").Parse("{{range .}} {{.}} {{end}}\n")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, items)
	if err != nil {
		panic(err)
	}
}

func BoolData() {
	tmpl, err := template.New("test").Parse("{{if .}}this is true{{else}}this is false{{end}}\n")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(os.Stdout, true)

}
func StructData() {
	user := User{
		Name: "John",
		Age:  30,
	}
	basic, err := template.New("struct").Parse("Hello, {{.Name}}! You are {{.Age}} years old.\n")
	if err != nil {
		panic(err)
	}
	basic.Execute(os.Stdout, user)
}
func Basic() {
	basic, err := template.New("basic").Parse("Hello, {{.}}!\n")
	if err != nil {
		panic(err)
	}
	basic.Execute(os.Stdout, "John")
}

// 模板变量使用示例
func TemplateVariables() {
	fmt.Println("\n=== 模板变量使用示例 ===")
	// 1. 使用 $ 变量
	tmpl1, _ := template.New("vars1").Parse(`
{{$name := "Alice"}}
{{$age := 250}}
Name: {{$name}}, Age: {{$age}}
	`)
	tmpl1.Execute(os.Stdout, nil)
}

// 复杂数据结构示例
func ComplexData() {
	fmt.Println("\n=== 复杂数据结构示例 ===")

	type Product struct {
		Name  string
		Price float64
		Stock int
	}

	type Store struct {
		Name     string
		Products []Product
	}

	store := Store{
		Name: "Tech Store",
		Products: []Product{
			{Name: "Laptop", Price: 999.99, Stock: 5},
			{Name: "Mouse", Price: 29.99, Stock: 20},
			{Name: "Keyboard", Price: 79.99, Stock: 15},
		},
	}

	tmpl, _ := template.New("complex").Parse(`
Store: {{.Name}}
{{range $index, $product := .Products}}
  {{$index | add1}}. {{$product.Name}}
     Price: ${{printf "%.2f" $product.Price}}
     Stock: {{$product.Stock}} units
     {{if lt $product.Stock 10}}⚠️ Low stock!{{end}}
{{end}}
`)

	tmpl = template.Must(tmpl.Funcs(template.FuncMap{
		"add1": func(i int) int { return i + 1 },
	}).Parse(tmpl.Tree.Root.String()))

	tmpl.Execute(os.Stdout, store)
}
