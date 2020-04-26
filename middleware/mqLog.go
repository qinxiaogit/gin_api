package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"time"
)

//记录日志到mq
func LoggerToMQ() gin.HandlerFunc {


	return func(c *gin.Context) {
		conn ,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err!= nil{
			log.Fatal("failed to connect to rabbitMQ:%v",err)
		}
		defer conn.Close()
		ch,err := conn.Channel()
		if err!= nil{
			log.Fatal("failed to connect to channel:",err)
		}
		q, err := ch.QueueDeclare("task_queue",
			true,
			false,
			false,
			false,
			nil,
		)
		if err!= nil{
			log.Fatal("failed to declare to queue:",err)
		}

		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		//执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		esLog := &MLog{Time: int64(latencyTime.Seconds()),
			Ip: c.ClientIP(), Status: c.Writer.Status(), Method: c.Request.Method, uri: c.Request.RequestURI}
		jsStr,_:=json.Marshal(esLog)
		err = ch.Publish("",q.Name,false,false,
			amqp.Publishing{
				ContentType:"application/json",
				Body:[]byte(jsStr),
				DeliveryMode:amqp.Persistent,//消息持久化
			},
				)
		if err!=nil{
			fmt.Errorf("push error",err)
		}
	}
}

//消息消费
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func MqConsumer(){
	conn,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch , err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable 队列持久化
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to open a channel")
	// 为了保证公平分发，不至于其中某个consumer一直处理，而其他不处理
	err = ch.Qos(
		1,     // prefetch count  在server收到consumer的ACK之前，预取的数量。为1，表示在没收到consumer的ACK之前，只会为其分发一个消息
		0,     // prefetch size 大于0时，表示在收到consumer确认消息之前，将size个字节保留在网络中
		false, // global  true:Qos对同一个connection的所有channel有效； false:Qos对同一个channel上的所有consumer有效
	)
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack   不进行自动ACK
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {  // msgs 是一个channel,从中取东西
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))  // 统计d.Body中的"."的个数
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second) // 有几个点就sleep几秒
			log.Printf("Done")
			d.Ack(false)  // 手动ACK，如果不ACK的话，那么无法保证这个消息被处理，可能它已经丢失了（比如消息队列挂了）
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever


}