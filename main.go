package main

import (
	"github.com/gin-gonic/gin"
	"go_gin_api/grpc/server/controller"
	"go_gin_api/grpc/server/proto/hello"
	servicListen "go_gin_api/grpc/server/proto/listen"
	"go_gin_api/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)
var wait = sync.WaitGroup{}
func main() {

	wait.Add(1)
	go runGrpc()
	r := gin.Default()
	//r.Use(middleware.LoggerToFile())
	//r.Use(middleware.LoggerToMongo())
	//r.Use(middleware.LoggerToES())
	//r.Use(middleware.LoggerToMQ())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.InitRouter(r)

	//go middleware.MqConsumer()
	r.Run()
	wait.Wait();
}

const (
	Address = "0.0.0.0:9090"
	)
func runGrpc(){
	defer wait.Done()
	listen ,err := net.Listen("tcp",Address)
	if err!= nil{
		log.Fatalf("Failed to listen: %v",err)
	}
	s:= grpc.NewServer()
	//服务注册
	hello.RegisterHelloServer(s,&controller.HelloController{})
	servicListen.RegisterListenServer(s,&controller.ListenController{})
	log.Println("listen on "+Address)

	grpc_health_v1.RegisterHealthServer(s,&controller.HealthImpl{})
	controller.RegisterToConsul()
	reflection.Register(s)
	if err := s.Serve(listen);err!=nil{
		log.Fatalf("Failed to serve: %v",err)
	}
}
