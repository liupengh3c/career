package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	jsoniter "github.com/json-iterator/go"
	"github.com/liupengh3c/esbuilder"
	"github.com/panjf2000/ants/v2"
)

// 最外层数据结构
type Documents struct {
	ScrollID string      `json:"_scroll_id"`
	Shards   Shards      `json:"_shards"`
	Hits     HitOutLayer `json:"hits"`
	TimedOut bool        `json:"timed_out"`
	Took     int         `json:"took"`
}
type Shards struct {
	Failed     int `json:"failed"`
	Skipped    int `json:"skipped"`
	Successful int `json:"successful"`
	Total      int `json:"total"`
}
type HitOutLayer struct {
	Hits     []Hits  `json:"hits"`
	MaxScore float64 `json:"max_score"`
	Total    Total   `json:"total"`
}
type Hits struct {
	ID     string         `json:"_id"`
	Index  string         `json:"_index"`
	Score  float64        `json:"_score"`
	Source map[string]any `json:"_source"`
	Type   string         `json:"_type"`
}
type Total struct {
	Relation string `json:"relation"`
	Value    int    `json:"value"`
}

func main() {
	var wg sync.WaitGroup

	date := flag.String("date", "World", "a name to say hello")
	threadCnt := flag.Int("thread_cnt", 8, "thread count")
	flag.Parse()
	tagIndex := "new_tag_" + *date
	threadpool, err := ants.NewPool(*threadCnt)
	if err != nil {
		fmt.Println("create thread pool err:", err.Error())
		return
	}
	client, err := NewEsClient()
	if err != nil {
		fmt.Println("create client err:", err.Error())
		return
	}
	fmt.Println("es connect success,begin to process "+tagIndex+" thread count:", *threadCnt)
	daysCount, err := DaysInMonthStr(*date)
	if err != nil {
		fmt.Println("get days count err:", err.Error())
		return
	}
	year, err := strconv.Atoi((*date)[:4])
	if err != nil {
		fmt.Println("invalid year:", err.Error())
		return
	}

	monthInt, err := strconv.Atoi((*date)[4:])
	if err != nil || monthInt < 1 || monthInt > 12 {
		fmt.Println("invalid month:", err.Error())
		return
	}
	fmt.Println(year, monthInt, daysCount, tagIndex)
	for i := 0; i < daysCount; i++ {
		st, et := DayTimestamps(year, monthInt, i+1)
		day := time.Date(year, time.Month(monthInt), i+1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		wg.Add(1)
		threadpool.Submit(func() {
			ScrollSearch(client, tagIndex, st, et, day, &wg)
		})
	}
	wg.Wait()
}
func NewEsClient() (*elasticsearch.Client, error) {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://xx.xx.xxx.xx:8200"},
		Username:  "your_es_username",
		Password:  "your_es_password",
		Transport: transport,
	})

	if err != nil {
		// fmt.Println("create client err:", err.Error())
		return client, err
	}
	return client, nil
}

func ScrollSearch(client *elasticsearch.Client, tagIndex string, st, et int64, day string, wg *sync.WaitGroup) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var scrollContextTime time.Duration = 25
	adbCount := 0
	adbProcessCnt := 0
	insertCount := 0
	var size int64 = 2000
	// overlapDocs := []Hits{}
	docs := Documents{}
	defer wg.Done()
	fmt.Println("begin to process day:", day)
	dslQuery := esbuilder.NewDsl()
	boolQuery := esbuilder.NewBoolQuery()
	boolQuery.Filter(esbuilder.NewTermQuery("tag_name.keyword", "adb"))
	boolQuery.Filter(esbuilder.NewRangeQuery("start_time").Gte(st).Lte(et))
	dslQuery.SetSize(size)
	dslQuery.SetQuery(boolQuery)
	dslQuery.SetOrder(esbuilder.NewSortQuery("start_time", "asc"))
	fmt.Println(dslQuery.BuildJson())
	begin := time.Now()
	resp, err := client.Search(
		client.Search.WithIndex(tagIndex),
		client.Search.WithBody(strings.NewReader(dslQuery.BuildJson())),
		client.Search.WithScroll(time.Minute*scrollContextTime),
	)
	if err != nil {
		fmt.Println("search err:", err.Error())
		return
	}
	fmt.Println("search time cost:", time.Since(begin).Seconds())
	err = json.NewDecoder(resp.Body).Decode(&docs)
	// 必须关闭
	resp.Body.Close()

	if err != nil {
		fmt.Println("decode err:", err.Error(), "day=", day)
		return
	}
	fmt.Println("all adb tags count:", docs.Hits.Total.Value, "date=", day)
	adbCount = docs.Hits.Total.Value
	begin = time.Now()
	for _, doc := range docs.Hits.Hits {
		insertCount += isHaveOverlapWithAdb(client, tagIndex, doc, day)
	}
	adbProcessCnt += len(docs.Hits.Hits)
	fmt.Println("search all other tags and BatchInSert time cost:", time.Since(begin).Seconds())
	if len(docs.Hits.Hits) < int(size) {
		fmt.Println("no more data to process,doc count:", len(docs.Hits.Hits), "size=", size)
		return
	}
	fmt.Printf("insert sucess count:%v, origin data info:%v/%v,date=%v\n", insertCount, adbProcessCnt, adbCount, day)
	scrollId := docs.ScrollID
	for {
		docs = Documents{}
		// overlapDocs = overlapDocs[:0]
		begin = time.Now()
		resp, err = client.Scroll(
			client.Scroll.WithScrollID(scrollId),
			client.Scroll.WithScroll(time.Minute*scrollContextTime),
		)
		if err != nil {
			fmt.Println("scroll err:", err.Error())
			return
		}
		json.NewDecoder(resp.Body).Decode(&docs)
		// 必须关闭
		resp.Body.Close()
		if len(docs.Hits.Hits) == 0 {
			fmt.Println("no more data to scroll,day=", day)
			break
		}
		adbProcessCnt += len(docs.Hits.Hits)
		fmt.Println("Scroll time cost:", time.Since(begin).Seconds(), "day=", day)
		begin = time.Now()
		// fmt.Println("search count:", len(docs.Hits.Hits))
		for _, doc := range docs.Hits.Hits {
			insertCount += isHaveOverlapWithAdb(client, tagIndex, doc, day)
		}

		fmt.Println("scroll all other tags and BatchInSert time cost:", time.Since(begin).Seconds(), "day=", day)
		fmt.Printf("scrool insert sucess count:%v, origin data info:%v/%v,date=%v\n", insertCount, adbProcessCnt, adbCount, day)
		scrollId = docs.ScrollID
	}

	client.ClearScroll(
		client.ClearScroll.WithScrollID(scrollId),
	)
}

func isHaveOverlapWithAdb(client *elasticsearch.Client, tagIndex string, doc Hits, day string) int {
	insertDocs := []Hits{}
	total := 0
	carId := doc.Source["car_id"].(string)
	st := int64(doc.Source["start_time"].(float64))
	et := int64(doc.Source["end_time"].(float64))
	allTags := []string{
		"pnc_point",
		"adc_scenario",
		"adc_behavior",
		"road_struc",
	}

	for _, tag := range allTags {
		docs := Documents{}
		dslQuery := esbuilder.NewDsl()
		boolQuery := esbuilder.NewBoolQuery()
		boolQuery.Filter(esbuilder.NewTermQuery("tag_name.keyword", tag))
		boolQuery.Filter(esbuilder.NewTermQuery("car_id.keyword", carId))
		boolQuery.Filter(esbuilder.NewRangeQuery("start_time").Lte(et))
		boolQuery.Filter(esbuilder.NewRangeQuery("end_time").Gte(st))
		dslQuery.SetQuery(boolQuery)
		// fmt.Println(dslQuery.BuildJson())
		resp, err := client.Search(
			client.Search.WithIndex(tagIndex),
			client.Search.WithBody(strings.NewReader(dslQuery.BuildJson())),
		)
		if err != nil {
			fmt.Println("search other tag err:", err.Error())
			return total
		}

		err = json.NewDecoder(resp.Body).Decode(&docs)
		// 必须关闭
		resp.Body.Close()
		if err != nil {
			fmt.Println("json decode err:", err.Error())
			continue
		}

		if len(docs.Hits.Hits) == 0 {
			// fmt.Println("tag_name:", tag, "car_id:", carId, "start_time:", st, "end_time:", et, "have no docs")
			continue
		}
		insertDocs = append(insertDocs, docs.Hits.Hits...)
	}
	if len(insertDocs) == 0 {
		return total
	}
	for i := 0; i < 3; i++ {
		err := BatchInSert(client, insertDocs)
		if err != nil {
			continue
		}
		fmt.Println("car_id:", carId, "start_time:", st, "end_time:", et, "insert success count:", len(insertDocs), "date=", day)
		total += len(insertDocs)
		break
	}
	return total
}
func BatchInSert(client *elasticsearch.Client, docs []Hits) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	buf := bytes.Buffer{}
	for _, doc := range docs {
		meta := map[string]any{
			"index": map[string]any{
				"_index": doc.Index + "_bak",
				"_id":    doc.ID,
				"_type":  "_doc",
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}
		if err := json.NewEncoder(&buf).Encode(doc.Source); err != nil {
			return err
		}
	}
	// resp, err := client.Bulk(&buf, client.Bulk.WithContext(context.Background()))
	req := esapi.BulkRequest{
		Body:   &buf,
		Pretty: true, // 格式化响应
	}
	resp, err := req.Do(context.Background(), client)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("error status code: ", resp.StatusCode)
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}
	return nil
}

// DaysInMonthStr 计算给定年月字符串（格式：YYYYMM）对应月份的天数
func DaysInMonthStr(ym string) (int, error) {
	if len(ym) != 6 {
		return 0, fmt.Errorf("invalid format: must be YYYYMM")
	}

	year, err := strconv.Atoi(ym[:4])
	if err != nil {
		return 0, fmt.Errorf("invalid year: %w", err)
	}

	monthInt, err := strconv.Atoi(ym[4:])
	if err != nil || monthInt < 1 || monthInt > 12 {
		return 0, fmt.Errorf("invalid month: %w", err)
	}

	month := time.Month(monthInt)

	// 构造下个月的第0天（即当前月最后一天）
	t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	return t.Day(), nil
}

func DayTimestamps(year int, month int, day int) (startTs, endTs int64) {
	// 构建当天 0 点
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	// 第二天 0 点 - 1 纳秒，就是当天的结束时间
	end := start.Add(24*time.Hour - time.Nanosecond)

	return start.UnixMilli(), end.UnixMilli()
}
