package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type DorisConfig struct {
	Host     string `json:"fe_host"`       // FE HTTP地址
	Port     int    `json:"fe_query_port"` // FE查询端口
	User     string `json:"user"`          // 用户名
	Password string `json:"password"`      // 密码
	Database string `json:"database"`      // 数据库
	Version  string `json:"version"`
}
type Record struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func StreamLoadJson(cfg DorisConfig, rows []Record) error {
	// 将 rows 转成 JSON Lines 格式
	var buf bytes.Buffer
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var isBig = false
	datas, _ := json.Marshal(rows)
	url := fmt.Sprintf(
		"http://%s:8030/api/%s/table_name/_stream_load",
		cfg.Host, cfg.Database,
	)
	if len(datas) > 104857600 {
		isBig = true
		for _, row := range rows {
			r, _ := json.MarshalToString(row)
			if len(r) > 1024*1024 {
				continue
			}
			buf.WriteString(r)
			buf.WriteString("\n")
		}
	} else {
		buf.Write(datas)
	}

	req, err := http.NewRequest("PUT", url, &buf)
	if err != nil {
		return err
	}

	// Header 配置
	label := fmt.Sprintf("streamload-%d", time.Now().UnixNano())
	req.Header.Set("label", label)
	req.Header.Set("format", "json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Expect", "100-continue")
	if isBig {
		req.Header.Set("read_json_by_line", "true")
		req.Header.Set("strip_outer_array", "false")
	} else {
		req.Header.Set("strip_outer_array", "true")
	}

	req.SetBasicAuth(cfg.User, cfg.Password)

	// HTTP 客户端
	client := &http.Client{
		Timeout: time.Minute * 2,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("stream load error:", err.Error())
		return fmt.Errorf("stream load failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Doris response:", string(body))

	if resp.StatusCode != 200 {
		return fmt.Errorf("stream load http code = %d", resp.StatusCode)
	}

	return nil
}

func main() {
	cfg := DorisConfig{}
	records := []Record{}
	StreamLoadJson(cfg, records)
}
