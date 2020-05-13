package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_gin_api/grpc/server/proto/hello"
	"google.golang.org/grpc"
	"io"
	"log"
)

const (
	Address = "127.0.0.1:9090"
)
func Grpc(c *gin.Context){
	//the secret sauce
	//resolver.SetDefaultScheme("dns")

	conn,err := grpc.Dial(Address,grpc.WithInsecure())
	if err != nil{
		log.Fatalln(err)
	}
	defer conn.Close()
	//初始化客户端
	client := hello.NewHelloClient(conn)
	//调用sayHello方法
	res , err := client.SayHello(context.Background(),&hello.HelloRequest{Name: "Hello World"})

	if err !=nil{
		log.Fatalln(err)
	}
	log.Println(res.Message)
	// 调用 LotsOfReplies 方法
	stream, err := client.LotsOfReplies(context.Background(),&hello.HelloRequest{Name: "Hello World replies"})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("stream.Recv: %v", err)
		}

		log.Printf("%s", res.Message)
	}
	c.JSON(200,gin.H{
		"v1" :"index",
		"name":"success",
		"price":200,
	})
}
