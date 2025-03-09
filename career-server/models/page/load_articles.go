package page

import (
	"career-server/lib/result"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

type Articles struct {
	BannerList []Banner    `json:"banner_list" toml:"banner_list"`
	GroupList  []GroupList `json:"group_list" toml:"group_list"`
}
type Banner struct {
	Img  string `json:"img" toml:"img"`
	Text string `json:"text" toml:"text"`
	URL  string `json:"url" toml:"url"`
}
type Items struct {
	ID    int    `json:"id" toml:"id"`
	Img   string `json:"img" toml:"img"`
	Title string `json:"title" toml:"title"`
	Desc  string `json:"desc" toml:"desc"`
	URL   string `json:"url" toml:"url"`
}
type GroupList struct {
	Name  string  `json:"name" toml:"name"`
	Show  bool    `json:"show" toml:"show"`
	Items []Items `json:"items" toml:"items"`
}

func LoadArticles(ctx *gin.Context, response *result.JsonResponseInterface) {
	articles := Articles{}
	// articles := map[string]interface{}{}
	curPath, _ := os.Getwd()
	fmt.Println("curPath:" + curPath)
	_, err := toml.DecodeFile(curPath+"/conf/articles.toml", &articles)
	if err != nil {
		fmt.Println("decode toml err:" + err.Error())
		return
	}
	fmt.Println(articles)
	response.Data = articles
}
