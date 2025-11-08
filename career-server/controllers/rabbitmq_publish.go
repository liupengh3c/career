package controllers

import (
	"career-server/resource"
	"fmt"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitMQPublish(ctx *gin.Context) {
	msg := `{"car_id":"ARCF082","start_time":1741598500863,"end_time":1741598509183,"tag_name":"pnc_point","tag_value":"njs_follow_scene","user":"liujianxin01","rule_id":1,"rule_version":"2.1.0.0","rule_priority":1,"update_time":1741685070,"tag_additional_info":{"driving_mode":1,"end_point":{"x":230210.716313,"y":3363380.046753,"z":9.312445},"obs_id":[],"rule_id":0,"source_type":"PLANNING","start_point":{"x":230330.474012,"y":3363391.718843,"z":10.375912}},"task_info":{"short_range_radar":"SRR208_21SSCL","task_purpose":6,"camera":"PANTHER","product":2,"map_region":"WuHanJingKaiQuKaiChengLuWang","hw_version":"HW5_0_R2","lidar_type":"PANDAR40P","hdmap_version":"hdmap-WuHanJingKaiQuKaiChengLuWang_1.14.93.327","iso_version":"6.4.4.16","long_range_radar":"ARS430","gnss":"BDN3_0"}}`
	err := resource.RabbitMQChannel.Publish(
		"tag_engine", // exchange
		"tag_insert", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
			// DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		fmt.Println("Failed to publish a message", err.Error())
		return
	}
}
