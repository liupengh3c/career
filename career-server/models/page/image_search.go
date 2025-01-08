package page

import (
	"bytes"
	"career-server/lib/result"
	"career-server/models/data"
	"career-server/resource"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func ImageSearch(ctx *gin.Context, fileName string, response *result.JsonResponseInterface) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	floatVec := make([]float64, 0)
	curPath, _ := os.Getwd()
	// 运行 Python 脚本并捕获输出
	cmd := exec.Command("python3", curPath+"/features/feature.py", curPath+"/images/")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running Python script:", err)
		return
	}

	// 打印 Python 脚本的输出
	vecs := out.String()
	start := strings.LastIndex(vecs, "[")
	end := strings.LastIndex(vecs, "]")
	content := vecs[start+1 : end-1]
	sliVecs := strings.Split(content, ", ")
	// 将字符串切片转换为浮点数切片,拿到特征向量
	for _, v := range sliVecs {
		valu, _ := strconv.ParseFloat(v, 64)
		floatVec = append(floatVec, valu)
	}
	mapQuery := map[string]any{
		"size": 10,
		"query": map[string]any{
			"knn": map[string]any{
				"field":        "IFV",
				"k":            10,
				"query_vector": floatVec,
			},
		},
	}
	dsl, _ := json.MarshalToString(mapQuery)
	results := data.ImageSearch(resource.ElasticClient, dsl)
	response.Data = results
}
