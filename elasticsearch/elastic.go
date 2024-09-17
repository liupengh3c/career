package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic/v7"
)

var condition = `{
    "task_info": {
        "time": [
            1725379200000,
            1725983999000
        ],
        "time_period": [
            "10-11",
            "11-12"
        ],
        "iso_version": [
            "7.2.27.205",
            "7.2.27.204"
        ],
        "map_region": [
            "hdmap-WuHanJingKaiQuKaiChengLuWang",
            "WuHanJingKaiQuKaiChengLuWang",
            "WuHanKaiChengCheLiangCeShiPaiZhaoKaoShi"
        ],
        "lidar_type": [
            "PANDAR90",
            "HESAI90",
            "AT128"
        ],
        "task_purpose": [
            6,
            4
        ],
        "hw_version": [
            "HW6_0_R2",
            "HW6_0_R1"
        ]
    },
    "tag_info": [
        {
            "primary_info": {
                "key": "DM_perception",
                "range": null
            },
            "additional_info": [
                {
                    "key": "obstacle_sub_type",
                    "values": [
                        "BigCar",
                        "SmallCar"
                    ],
                    "range": null
                },
                {
                    "key": "perception_problem",
                    "values": [
                        "Boundary_Error",
                        "Speed_delay"
                    ],
                    "range": null
                }
            ]
        },
        {
            "primary_info": {
                "key": "adc_behavior",
                "values": [
                    "LEFT_TURN",
                    "RIGHT_TURN"
                ],
                "range": null
            },
            "additional_info": [
                {
                    "key": "adc_driving_mode",
                    "values": [
                        "COMPLETE_AUTO_DRIVE"
                    ],
                    "range": null
                },
                {
                    "key": "average_speed",
                    "range": {
                        "gte": 10,
                        "lte": 20
                    }
                }
            ]
        },
        {
            "primary_info": {
                "key": "e2e_adc_behavior",
                "values": [
                    "front_vehicle_brake"
                ],
                "range": null
            },
            "additional_info": null
        },
        {
            "primary_info": {
                "key": "adc_scenario",
                "values": [
                    "Highway"
                ],
                "range": null
            },
            "additional_info": [
                {
                    "key": "static_scene",
                    "values": [
                        "JunctionGostraight",
                        "JunctionLeftTurn"
                    ],
                    "range": null
                },
                {
                    "key": "dynamic_scene",
                    "values": [
                        "ChangingLane"
                    ],
                    "range": null
                }
            ]
        },
        {
            "primary_info": {
                "key": "grading_entry",
                "values": [
                    "brake"
                ],
                "range": null
            },
            "additional_info": null
        },
        {
            "primary_info": {
                "key": "DM_scene",
                "range": null
            },
            "additional_info": [
                {
                    "key": "scene_label",
                    "values": [
                        "Obs_By_BigCar"
                    ],
                    "range": null
                }
            ]
        },
        {
            "primary_info": {
                "key": "weather",
                "values": [
                    "SUNY_MODE",
                    "RAIN_MODE",
                    "FOG_MODE"
                ],
                "range": null
            },
            "additional_info": null
        },
        {
            "primary_info": {
                "key": "perception_mining_rough",
                "range": null
            },
            "additional_info": [
                {
                    "key": "dm_rule_id",
                    "values": [
                        "82387",
                        "9630"
                    ],
                    "range": null
                },
                {
                    "key": "mining_task_id",
                    "values": [
                        "2123"
                    ],
                    "range": null
                },
                {
                    "key": "case_tag",
                    "values": [
                        "FALSE_DETECT_VRU"
                    ],
                    "range": null
                },
                {
                    "key": "obstacle_type",
                    "values": [
                        "BICYCLE"
                    ],
                    "range": null
                },
                {
                    "key": "obstacle_sub_type",
                    "values": [
                        "BUS"
                    ],
                    "range": null
                },
                {
                    "key": "light_type",
                    "values": [
                        "EMERGENCY_LIGHT",
                        "RIGHT_TURN_LIGHT"
                    ],
                    "range": null
                },
                {
                    "key": "scene_type",
                    "values": [
                        "HIGHWAY",
                        "URBAN"
                    ],
                    "range": null
                },
                {
                    "key": "reversing_type",
                    "values": [
                        "OPPOSITE_3POINT"
                    ],
                    "range": null
                }
            ]
        },
        {
            "primary_info": {
                "key": "self_define_rule",
                "range": null
            },
            "additional_info": [
                {
                    "key": "mining_job_id",
                    "values": [
                        "114"
                    ],
                    "range": null
                }
            ]
        },
        {
            "primary_info": {
                "key": "obs_feature",
                "range": null
            },
            "additional_info": [
                {
                    "key": "type",
                    "values": [
                        "CAR",
                        "VAN"
                    ],
                    "range": null
                },
                {
                    "key": "min_dist",
                    "range": {
                        "gte": 1,
                        "lte": 22
                    }
                },
                {
                    "key": "width",
                    "range": {
                        "gte": 1,
                        "lte": 3
                    }
                },
                {
                    "key": "length",
                    "range": {
                        "gte": 2,
                        "lte": 2
                    }
                },
                {
                    "key": "average_speed",
                    "range": {
                        "gte": 1,
                        "lte": 4
                    }
                }
            ]
        },
        {
            "primary_info": {
                "key": "pnc_point",
                "values": [
                    "cyclist_predecision_handle_points_by_cross_map_curb"
                ],
                "range": null
            },
            "additional_info": [
                {
                    "key": "creator",
                    "values": [
                        "offline"
                    ],
                    "range": null
                }
            ]
        }
    ],
    "postprocess_settings": {

    }
}`

type SearchCondition struct {
	TaskInfo            TaskInfo            `json:"task_info"`
	TagInfo             []TagInfo           `json:"tag_info"`
	PostprocessSettings PostprocessSettings `json:"postprocess_settings"`
}
type TaskInfo struct {
	Time        []int64  `json:"time"`
	TimePeriod  []string `json:"time_period"`
	IsoVersion  []string `json:"iso_version"`
	MapRegion   []string `json:"map_region"`
	LidarType   []string `json:"lidar_type"`
	TaskPurpose []int    `json:"task_purpose"`
	HwVersion   []string `json:"hw_version"`
}

type AdditionalInfo struct {
	Key    string      `json:"key"`
	Values []string    `json:"values"`
	Range  interface{} `json:"range"`
}
type PrimaryInfo struct {
	Key    string      `json:"key"`
	Values []string    `json:"values"`
	Range  interface{} `json:"range"`
}

type TagInfo struct {
	PrimaryInfo    PrimaryInfo      `json:"primary_info,omitempty"`
	AdditionalInfo []AdditionalInfo `json:"additional_info"`
}
type PostprocessSettings struct {
}
type EsInstance struct {
	Client *elastic.Client    // es客户端
	Index  map[string]Indexes // 所有索引
}

type Indexes struct {
	Name    string `json:"name"`    // 索引名称
	Mapping string `json:"mapping"` // 索引结构
}

func NewEsInstance() (*EsInstance, error) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"), // 支持多个服务地址，逗号分隔
		elastic.SetBasicAuth("user_name", "user_password"),
		elastic.SetSniff(false), // 跳过ip检查，默认是true
	)
	if err != nil {
		return &EsInstance{}, err
	}
	return &EsInstance{
		Client: client,
		Index:  map[string]Indexes{},
	}, nil
}

func main() {
	var searchConditions SearchCondition
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.UnmarshalFromString(condition, &searchConditions)
	for _, condition := range searchConditions.TagInfo {
		boolQuery := elastic.NewBoolQuery()
		boolQuery.Filter(elastic.NewTermQuery("tag_name.keyword", condition.PrimaryInfo.Key))
		if len(condition.PrimaryInfo.Values) > 0 {
			tt := []interface{}{}
			for _, v := range condition.PrimaryInfo.Values {
				tt = append(tt, v)
			}
			boolQuery.Filter(elastic.NewTermsQuery(condition.PrimaryInfo.Key, tt...))
		}
		s, _ := boolQuery.Source()
		st, _ := json.MarshalToString(s)
		fmt.Println(st)
	}
}
