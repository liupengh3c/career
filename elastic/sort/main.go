package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/panjf2000/ants/v2"

	jsoniter "github.com/json-iterator/go"
	"github.com/liupengh3c/esbuilder"
)

// 最外层数据结构
type Documents struct {
	ScrollID     string       `json:"_scroll_id"`
	Shards       Shards       `json:"_shards"`
	Hits         HitOutLayer  `json:"hits"`
	TimedOut     bool         `json:"timed_out"`
	Took         int          `json:"took"`
	Aggregations Aggregations `json:"aggregations"`
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
type Buckets struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
}
type Aggs struct {
	DocCountErrorUpperBound int       `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int       `json:"sum_other_doc_count"`
	Buckets                 []Buckets `json:"buckets"`
}
type Aggregations struct {
	Aggs Aggs `json:"aggs"`
}
type AdcDriverStatus struct {
	Status    int     `json:"status"`
	StartTime float64 `json:"start_time"`
	Miles     float64 `json:"miles"`
	EndTime   float64 `json:"end_time"`
	Duration  float64 `json:"duration"`
}
type ByStartTime []AdcDriverStatus

func (a ByStartTime) Len() int {
	return len(a)
}
func (a ByStartTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByStartTime) Less(i, j int) bool {
	return a[i].StartTime < a[j].StartTime
}

type AdcDriverStatusEs struct {
	CarID     string `json:"car_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	// TaskInfo          TaskInfo          `json:"task_info"`
	DwTagVersion      string            `json:"dw_tag_version"`
	CreateTime        string            `json:"create_time"`
	CreateUser        string            `json:"create_user"`
	TagSource         string            `json:"tag_source"`
	UpdateTime        time.Time         `json:"update_time"`
	TagAdditionalInfo TagAdditionalInfo `json:"tag_additional_info"`
	TagName           string            `json:"tag_name"`
	AdcDrivingStatus  string            `json:"adc_driving_status"`
}
type TaskInfo struct {
	TaskPurpose int `json:"task_purpose"`
}
type TagAdditionalInfo struct {
	DumpDate       string `json:"dump_date"`
	AdcDrivingMode string `json:"adc_driving_mode"`
	TaskID         string `json:"task_id"`
}

// var taskPurpose = map[string]int{
// 	"debug":         0,
// 	"ads":           1,
// 	"test":          2,
// 	"collection":    3,
// 	"dailybuild":    4,
// 	"roadtest":      5,
// 	"operation":     6,
// 	"mapcollection": 7,
// 	"prerelease":    8,
// 	"prepublish":    9,
// 	"publish":       10,
// 	"mapchecking":   11,
// 	"stationtest":   12,
// }

func NewJsEsClient() (*elasticsearch.Client, error) {
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
		Addresses: []string{"http://10.27.139.114:8200"},
		Username:  "superuser",
		Password:  "adudata@123",
		Transport: transport,
	})

	if err != nil {
		// fmt.Println("create client err:", err.Error())
		return client, err
	}
	return client, nil
}
func NewDwEsClient() (*elasticsearch.Client, error) {
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
		Addresses: []string{"http://10.21.238.15:8200"},
		Username:  "superuser",
		Password:  "ES_admin",
		Transport: transport,
	})

	if err != nil {
		// fmt.Println("create client err:", err.Error())
		return client, err
	}
	return client, nil
}
func main() {
	var wg sync.WaitGroup
	date := flag.String("date", "World", "the date for index")
	day := flag.String("day", "", "the date for index")
	threadCnt := flag.Int("thread_cnt", 8, "thread count")
	flag.Parse()

	threadpool, err := ants.NewPool(*threadCnt)
	if err != nil {
		fmt.Println("create thread pool err:", err.Error())
		return
	}
	jsClient, err := NewJsEsClient()
	if err != nil {
		fmt.Println("create js client err:", err.Error())
		return
	}
	dwClient, err := NewDwEsClient()
	if err != nil {
		fmt.Println("create dw client err:", err.Error())
		return
	}

	if *day != "" {
		if *day == "default" {
			one := time.Now().Add(-48 * time.Hour).Format("20060102")
			tagIndex := "auto_car.app_data_" + one[:6]
			fmt.Println("es connect success,begin to process "+tagIndex+" thread count:", *threadCnt)
			st, et := GetOneDayTimestamp(one)
			fmt.Println("only process one day:", one, "st", st, "et", et)
			tasks := SearchAllTasks(jsClient, tagIndex, st, et, one)
			fmt.Println("begin to process day:", one, "task_id count:", len(tasks))
			for _, task := range tasks {
				wg.Add(1)
				threadpool.Submit(func() {
					SearchAdcDriverStatusByTaskId(jsClient, dwClient, tagIndex, task, one, &wg)
				})
			}
		} else {
			tagIndex := "auto_car.app_data_" + (*day)[:6]
			fmt.Println("es connect success,begin to process "+tagIndex+" thread count:", *threadCnt)
			st, et := GetOneDayTimestamp(*day)
			tasks := SearchAllTasks(jsClient, tagIndex, st, et, *day)
			fmt.Println("begin to process day:", *day, "task_id count:", len(tasks))
			for _, task := range tasks {
				wg.Add(1)
				threadpool.Submit(func() {
					SearchAdcDriverStatusByTaskId(jsClient, dwClient, tagIndex, task, *day, &wg)
				})
			}
		}
	} else {
		tagIndex := "auto_car.app_data_" + *date
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
		fmt.Printf("year: %d, month: %d, days count: %d, tag index: %s\n", year, monthInt, daysCount, tagIndex)
		for i := 0; i < daysCount; i++ {
			st, et := DayTimestamps(year, monthInt, i+1)
			day := time.Date(year, time.Month(monthInt), i+1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
			tasks := SearchAllTasks(jsClient, tagIndex, st, et, day)
			fmt.Println("begin to process day:", day, "task_id count:", len(tasks))
			for _, task := range tasks {
				wg.Add(1)
				threadpool.Submit(func() {
					SearchAdcDriverStatusByTaskId(jsClient, dwClient, tagIndex, task, day, &wg)
				})
			}
		}
	}
	// SearchAdcDriverStatusByTaskId(client, tagIndex, "ARCF1002_20250610140000")
	wg.Wait()
}

func SearchAllTasks(client *elasticsearch.Client, tagIndex string, st, et int64, day string) []string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	taskids := []string{}
	docs := Documents{}
	fmt.Println("begin to process day:", day)
	dslQuery := esbuilder.NewDsl()
	boolQuery := esbuilder.NewBoolQuery()
	boolQuery.Filter(esbuilder.NewTermQuery("data_type.keyword", "drive_status"))
	boolQuery.Filter(esbuilder.NewRangeQuery("start_time").Gte(fmt.Sprintf("%.2f", float32(st)/1000.0)))
	boolQuery.Filter(esbuilder.NewRangeQuery("start_time").Lte(fmt.Sprintf("%.2f", float32(et)/1000.0)))
	dslQuery.SetSize(0)
	dslQuery.SetQuery(boolQuery)
	dslQuery.SetOrder(esbuilder.NewSortQuery("start_time.keyword", "asc"))

	aggsQuery := esbuilder.NewAggsQuery("aggs")
	aggsQuery.Terms(esbuilder.NewAggsTerm("task_id.keyword", 10000))
	dslQuery.SetAggs(aggsQuery)
	fmt.Println(dslQuery.BuildJson())
	// begin := time.Now()
	resp, err := client.Search(
		client.Search.WithIndex(tagIndex),
		client.Search.WithBody(strings.NewReader(dslQuery.BuildJson())),
	)
	if err != nil {
		fmt.Println("search err:", err.Error())
		return taskids
	}
	// fmt.Println("search time cost:", time.Since(begin).Seconds())
	err = json.NewDecoder(resp.Body).Decode(&docs)
	// 必须关闭
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("decode err:", err.Error(), "day=", day)
		return taskids
	}
	for _, v := range docs.Aggregations.Aggs.Buckets {
		taskids = append(taskids, v.Key)
	}
	return taskids
}

func SearchAdcDriverStatusByTaskId(client *elasticsearch.Client, dwClient *elasticsearch.Client, tagIndex string, taskid, day string, wg *sync.WaitGroup) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	allTags := []AdcDriverStatusEs{}
	var size int64 = 200
	auto := []AdcDriverStatus{}
	manu := []AdcDriverStatus{}
	safe := []AdcDriverStatus{}
	docs := Documents{}
	defer wg.Done()
	dslQuery := esbuilder.NewDsl()
	boolQuery := esbuilder.NewBoolQuery()
	boolQuery.Filter(esbuilder.NewTermQuery("data_type.keyword", "drive_status"))
	boolQuery.Filter(esbuilder.NewTermQuery("task_id.keyword", taskid))
	dslQuery.SetSize(size)
	dslQuery.SetQuery(boolQuery)
	dslQuery.SetOrder(esbuilder.NewSortQuery("start_time.keyword", "asc"))

	// begin := time.Now()
	resp, err := client.Search(
		client.Search.WithIndex(tagIndex),
		client.Search.WithBody(strings.NewReader(dslQuery.BuildJson())),
	)
	if err != nil {
		fmt.Println("search err:", err.Error())
		return
	}
	// fmt.Println("search time cost:", time.Since(begin).Seconds())
	err = json.NewDecoder(resp.Body).Decode(&docs)
	// 必须关闭
	defer resp.Body.Close()

	if err != nil {
		return
	}

	if len(docs.Hits.Hits) == 0 {
		return
	}
	for _, doc := range docs.Hits.Hits {
		a, m, s := ProcessDriveData(doc.Source["data"].(string))
		auto = append(auto, a...)
		manu = append(manu, m...)
		safe = append(safe, s...)
	}
	allTags = append(allTags, DataToTag(docs.Hits.Hits[0], auto, "COMPLETE_AUTO_DRIVE")...)
	allTags = append(allTags, DataToTag(docs.Hits.Hits[0], manu, "COMPLETE_MANUAL")...)
	allTags = append(allTags, DataToTag(docs.Hits.Hits[0], safe, "SECURITY_MODE")...)
	if len(allTags) == 0 {
		fmt.Println("no valid data:", taskid)
		return
	}
	fmt.Println("the adc driving status tag count:", len(allTags), "taskid=", taskid, "day:", day)
	// for _, v := range allTags {
	// 	fmt.Println(time.UnixMilli(v.StartTime).Format("2006-01-02 15:04:05"), time.UnixMilli(v.EndTime).Format("2006-01-02 15:04:05"))
	// }

	for i := 0; i < 3; i++ {
		err := BatchInSert(dwClient, allTags)
		if err != nil {
			fmt.Println("insert es err:", err.Error())
			continue
		}
		fmt.Println("task_id:", taskid, "insert success docs count:", len(allTags))
		break
	}
}
func DataToTag(hit Hits, adcStatus []AdcDriverStatus, status string) []AdcDriverStatusEs {
	adcTags := []AdcDriverStatusEs{}
	sort.Sort(ByStartTime(adcStatus))
	if len(adcStatus) == 0 {
		return adcTags
	}
	current := adcStatus[0]
	for _, v := range adcStatus[1:] {
		if int64(current.EndTime*1000+1000) >= int64(v.StartTime*1000) {
			current.EndTime = v.EndTime
		} else {
			autoTag := AdcDriverStatusEs{
				CarID:     hit.Source["car_id"].(string),
				StartTime: int64(current.StartTime * 1000),
				EndTime:   int64(current.EndTime * 1000),
				// TaskInfo: TaskInfo{
				// 	TaskPurpose: taskPurpose[hit.Source["task_purpose"].(string)],
				// },
				DwTagVersion: "dw_2.0",
				CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
				CreateUser:   "liupeng17",
				TagSource:    "script",
				UpdateTime:   time.Now(),
				TagAdditionalInfo: TagAdditionalInfo{
					DumpDate:       fmt.Sprintf("%v", int64(hit.Source["date"].(float64))),
					AdcDrivingMode: status,
					TaskID:         hit.Source["task_id"].(string),
				},
				TagName:          "adc_driving_status",
				AdcDrivingStatus: status,
			}
			adcTags = append(adcTags, autoTag)
			current = v
		}
	}
	autoTag := AdcDriverStatusEs{
		CarID:     hit.Source["car_id"].(string),
		StartTime: int64(current.StartTime * 1000),
		EndTime:   int64(current.EndTime * 1000),
		// TaskInfo: TaskInfo{
		// 	TaskPurpose: taskPurpose[hit.Source["task_purpose"].(string)],
		// },
		DwTagVersion: "dw_2.0",
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		CreateUser:   "liupeng17",
		TagSource:    "script-ipipe",
		UpdateTime:   time.Now(),
		TagAdditionalInfo: TagAdditionalInfo{
			DumpDate:       fmt.Sprintf("%v", int64(hit.Source["date"].(float64))),
			AdcDrivingMode: status,
			TaskID:         hit.Source["task_id"].(string),
		},
		TagName:          "adc_driving_status",
		AdcDrivingStatus: status,
	}
	adcTags = append(adcTags, autoTag)
	return adcTags
}

// ProcessDriveData 处理从ADC驱动器接收到的数据，并返回处理后的结果
//
// 参数：
// data string - 从ADC驱动器接收到的原始数据，格式为JSON字符串
//
// 返回值：
// []map[string]int64 - 自动驾驶数据，
// []map[string]int64 - 人工驾驶数据，
// []map[string]int64 - 安全模式
func ProcessDriveData(data string) ([]AdcDriverStatus, []AdcDriverStatus, []AdcDriverStatus) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var adcDriverStatus []AdcDriverStatus
	adcAutoDriverResults := []AdcDriverStatus{}
	adcManuDriverResults := []AdcDriverStatus{}
	adcSafeDriverResults := []AdcDriverStatus{}
	err := json.Unmarshal([]byte(data), &adcDriverStatus)
	if err != nil {
		fmt.Println("Unmarshal err:", err.Error())
		return adcAutoDriverResults, adcManuDriverResults, adcSafeDriverResults
	}

	// 自动驾驶
	for _, v := range adcDriverStatus {
		if v.Status == 1 {
			one := AdcDriverStatus{}
			one.StartTime = v.StartTime
			one.EndTime = v.EndTime
			adcAutoDriverResults = append(adcAutoDriverResults, one)
		}
	}
	// 人工驾驶
	for _, v := range adcDriverStatus {
		if v.Status == 0 {
			one := AdcDriverStatus{}
			one.StartTime = v.StartTime
			one.EndTime = v.EndTime
			adcManuDriverResults = append(adcManuDriverResults, one)
		}
	}
	// 安全模式
	for _, v := range adcDriverStatus {
		if v.Status == 4 {
			one := AdcDriverStatus{}
			one.StartTime = v.StartTime
			one.EndTime = v.EndTime
			adcSafeDriverResults = append(adcSafeDriverResults, one)
		}
	}
	return adcAutoDriverResults, adcManuDriverResults, adcSafeDriverResults
}
func BatchInSert(client *elasticsearch.Client, docs []AdcDriverStatusEs) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	buf := bytes.Buffer{}
	for _, doc := range docs {
		index := "new_tag_" + time.UnixMilli(doc.StartTime).Format("200601")
		id := doc.CarID + "_" + doc.TagName + "_" + doc.AdcDrivingStatus + "_" + fmt.Sprintf("%v", doc.StartTime) + "_" + fmt.Sprintf("%v", doc.EndTime)
		meta := map[string]any{
			"index": map[string]any{
				"_index": index,
				"_id":    id,
				"_type":  "_doc",
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
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
		fmt.Println("error status code: ", resp)
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

func GetOneDayTimestamp(day string) (int64, int64) {
	// 解析日期字符串为 time.Time 对象（以北京时间为例）
	loc, _ := time.LoadLocation("Asia/Shanghai")
	dayStart, err := time.ParseInLocation("20060102", day, loc)
	if err != nil {
		panic(err)
	}

	// 当天起始时间
	startMillis := dayStart.UnixNano() / 1e6

	// 当天结束时间（即第二天 0 点减 1 毫秒）
	dayEnd := dayStart.Add(24 * time.Hour).Add(-time.Millisecond)
	endMillis := dayEnd.UnixNano() / 1e6

	fmt.Println("Start timestamp (ms):", startMillis)
	fmt.Println("End timestamp (ms):  ", endMillis)
	return startMillis, endMillis
}
