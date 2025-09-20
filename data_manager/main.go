package datamanager

const (
	INTERSECTION_LEFT = iota
	INTERSECTION_OUTER
	INTERSECTION_INNER
	INTERSECTION_RIGHT
)

type DataManager struct {
}

type DataMarker struct {
	TaskId      string `json:"task_id"`
	CarId       string `json:"car_id"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	Topics      string `json:"topics"`     // topic list
	ExpiredAt   int64  `json:"expired_at"` // expire time, unit: second,固定大于end_time 1个月或者2个月或其他整数月
	ExpireHuman string `json:"expire_human"`
}
