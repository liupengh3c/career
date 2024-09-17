package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MiniprogramOrder struct {
	gorm.Model                 //gorm预定义结构体
	Uid                string  `json:"uid"`
	OrderSn            string  `json:"order_sn"`
	Name               string  `json:"name"`
	Sex                string  `json:"sex"`
	IdentifyCardNumber string  `json:"identify_card_number"`
	Country            string  `json:"country"`
	IdentifyCardType   string  `json:"identify_card_type"`
	BornDate           string  `json:"born_date"`
	Nation             string  `json:"nation"`
	PhoneNum           string  `json:"phone_number"`
	Email              string  `json:"email"`
	Address            string  `json:"address"`
	CommodityId        string  `json:"commodity_id"`
	PlayType           string  `json:"play_type"`
	ApplicationPay     float32 `json:"application_pay"`
	Status             int     `json:"status"`
}

func initMysql() *gorm.DB {
	var err error
	dsn := "leguan:leguan@123@tcp(127.0.0.1:3306)/leguan?charset=utf8mb4&parseTime=True&loc=Local"
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database" + err.Error())
		return nil
	}
	fmt.Println("db connect success!")
	return client
}

// 根据uid查询订单
func Select(db *gorm.DB, uid string) []MiniprogramOrder {
	orderList := []MiniprogramOrder{}
	db.Where("uid = ?", uid).Find(&orderList)
	return orderList
}
func main() {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	db := initMysql()
	if db == nil {
		return
	}
	orderList := Select(db, "liupeng20240908")
	str, _ := json.MarshalToString(orderList)
	fmt.Println(str)
}
