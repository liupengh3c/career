package html

import (
	"net/http"

	"golang.org/x/net/html"
)

func DecodeHtml() {
	// 发送HTTP GET请求
	resp, err := http.Get("https://www.douyin.com/video/7354021255869156627")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return
	// }

	// 解析HTML文档
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	// 遍历文档并找到特定的元素
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			// for _, a := range n.Attr {
			// 	if a.Key == "id" && a.Val == "RENDER_DATA" {
			// 		fmt.Println(n.FirstChild.Data) // 打印script标签的文本内容
			// 	}
			// }
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}
