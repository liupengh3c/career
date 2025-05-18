package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ExcelRead(name string) [][]string {
	contents := [][]string{}
	f, err := excelize.OpenFile(name)
	if err != nil {
		fmt.Println("open file error:", err)
		return contents
	}
	defer f.Close()
	// 获取所有sheet的名称
	sheets := f.GetSheetMap()
	for _, name := range sheets {
		// 获取单个sheet的全部内容，正式上线需要考虑内存占用问题
		rows, err := f.GetRows(name)
		if err != nil {
			fmt.Println("get rows error:", err)
			return contents
		}
		fmt.Println("rows lenth:", len(rows))

		// 每次读取一行，减少内存占用
		for i := 0; i < 200; i++ {
			st, err := f.GetCellValue(name, fmt.Sprintf("A%v", i+2))
			if err != nil {
				fmt.Println("get cell value error:", err)
			}
			if st == "" {
				break
			}
			et, err := f.GetCellValue(name, fmt.Sprintf("B%v", i+2))
			if err != nil {
				fmt.Println("get cell value error:", err)
			}
			sub, err := f.GetCellValue(name, fmt.Sprintf("C%v", i+2))
			if err != nil {
				fmt.Println("get cell value error:", err)
			}
			fmt.Println(st, et, sub)
		}
		contents = append(contents, rows[1:]...)
	}
	return contents
}

func ExcelWrite(sheet1 [][]string, sheet2 [][]string) {
	col := []string{"A", "B", "C"}
	f := excelize.NewFile()
	defer f.Close()
	// 创建名为车端挖掘数据的工作表
	_, err := f.NewSheet("车端挖掘数据")
	if err != nil {
		fmt.Println("new sheet error:", err)
		return
	}
	// 设置单元格的值。
	f.SetCellValue("车端挖掘数据", "A1", "start_time")
	f.SetCellValue("车端挖掘数据", "B1", "end_time")
	f.SetCellValue("车端挖掘数据", "C1", "subject")
	for i, v := range sheet1 {
		for k, v1 := range v {
			f.SetCellValue("车端挖掘数据", col[k]+fmt.Sprintf("%v", i+2), v1)
		}
	}

	// 创建名为车端挖掘数据的工作表
	sheetName := "云端挖掘数据"
	index, err := f.NewSheet("云端挖掘数据")
	if err != nil {
		fmt.Println("new sheet error:", err)
		return
	}
	// 设置默认sheet，打开表格时，默认显示的sheet
	f.SetActiveSheet(index)
	// 设置单元格的值。
	f.SetCellValue(sheetName, "A1", "start_time")
	f.SetCellValue(sheetName, "B1", "end_time")
	f.SetCellValue(sheetName, "C1", "subject")
	for i, v := range sheet1 {
		for k, v1 := range v {
			f.SetCellValue(sheetName, col[k]+fmt.Sprintf("%v", i+2), v1)
		}
	}
	// 将Excel另存为文件
	if err := f.SaveAs("/Users/liupeng/Downloads/test.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func CreateChart() {
	f := excelize.NewFile()
	defer f.Close()

	contents := [][]any{
		{nil, "Apple", "Orange", "Pear"},
		{"Small", 2, 3, 3},
		{"Normal", 5, 2, 4},
		{"Large", 6, 7, 8},
	}
	for idx, row := range contents {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetSheetRow("Sheet1", cell, &row)
	}
	if err := f.AddChart("Sheet1", "E1", &excelize.Chart{
		Type: excelize.Col3DClustered,
		Series: []excelize.ChartSeries{
			{
				Name:       "Sheet1!$A$2",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$2:$D$2",
			},
			{
				Name:       "Sheet1!$A$3",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$3:$D$3",
			},
			{
				Name:       "Sheet1!$A$4",
				Categories: "Sheet1!$B$1:$D$1",
				Values:     "Sheet1!$B$4:$D$4",
			}},
		Title: []excelize.RichTextRun{
			{
				Text: "Fruit Chart",
			},
		},
	}); err != nil {
		fmt.Println(err)
		return
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs("/Users/liupeng/Downloads/Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
func main() {
	c := ExcelRead("/Users/liupeng/Downloads/final_new_0301.xlsx")
	ExcelWrite(c, c)
	// CreateChart()
}
