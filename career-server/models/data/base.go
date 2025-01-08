package data

// 最外层数据结构
type Documents struct {
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
	Sort   []any          `json:"sort"`
}
type Total struct {
	Relation string `json:"relation"`
	Value    int    `json:"value"`
}
