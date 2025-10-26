package resource

import (
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

var (
	ElasticClient *elasticsearch.Client

	Mux          sync.Mutex
	IpLimiter    map[string]*rate.Limiter
	IpLimiterCnt = 2
	IpLimiterMax = 5

	GlobalLimiter    *rate.Limiter
	GlobalLimiterCnt = 10
	GlobalLimiterMax = 20
	RedisClient      *redis.Client
)
