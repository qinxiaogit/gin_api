package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_gin_api/common"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MLog struct {
	Status int
	Time   int64
	Ip     string
	Method string
	uri    string
}

func LoggerToFile() gin.HandlerFunc {
	//日志文件
	fileName := path.Join(common.Log_FILE_PATH, common.Log_FILE_NAME)
	fmt.Println(fileName)
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
		panic(err)
	}

	//实例化日志
	logger := logrus.New()
	//日志目标
	logger.Out = src
	logger.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
	//日志级别
	logger.SetLevel(logrus.DebugLevel)
	//日志格式
	logger.SetFormatter(&logrus.TextFormatter{})
	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//fmt.Println("处理请求前:  ")
		//处理请求
		c.Next()
		//fmt.Println("处理请求后:  ")
		//结束时间
		endTime := time.Now()
		//执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		logger.Infof("| %3d| %13v |%15s|%s|",
			c.Writer.Status(),
			latencyTime,
			c.ClientIP(),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}

//记录日志到mongo
func LoggerToMongo() gin.HandlerFunc {
	//set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	//connect to mongoDb
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	//check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to mongoDb")
	collection := client.Database("test").Collection("log")

	return func(c *gin.Context) {
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		//执行时间
		latencyTime := endTime.Sub(startTime)
		ret, _ := collection.InsertOne(context.TODO(), MLog{Time: int64(latencyTime.Seconds()), Ip: c.ClientIP(), Status: c.Writer.Status(), Method: c.Request.Method, uri: c.Request.RequestURI})
		fmt.Println(ret)
	}
}

//记录日志到es
func LoggerToES() gin.HandlerFunc {
	var (
		r map[string]interface{}
		//wg sync.WaitGroup
	)
	//初始化es客户端
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client :%s", err)
	}
	//get Cluster info
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response :%s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Fatalf("ERROR getting response:%s", err)
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	var indexName = "go_gin_api"
	//set up the request object
	 i  := 0
	return func(c *gin.Context) {
		i++
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		//执行时间
		latencyTime := endTime.Sub(startTime)
		esLog := &MLog{Time: int64(latencyTime.Seconds()),
			Ip: c.ClientIP(), Status: c.Writer.Status(), Method: c.Request.Method, uri: c.Request.RequestURI}
		jsStr,_:=json.Marshal(esLog)
		//fmt.Println(string(jsStr))
		req  := esapi.IndexRequest{
			Index:indexName,
			DocumentID:strconv.Itoa(i),
			Body: strings.NewReader(string(jsStr)),
			Refresh:"true",
		}
		res ,err := req.Do(context.Background(),es)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
		} else {
			// Deserialize the response into a map.
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Printf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and indexed document version.
				log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			}
		}
	}
}


