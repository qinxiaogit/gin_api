package controller

import (
	"context"
	"fmt"
	"go_gin_api/grpc/plugs/consul"
	"go_gin_api/grpc/server/proto/hello"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type HelloController struct{}

func (h *HelloController) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{Message: fmt.Sprintf("%s", in.Name)}, nil
}

func (h *HelloController) LotsOfReplies(in *hello.HelloRequest, stream hello.Hello_LotsOfRepliesServer) error {
	for i := 0; i < 10; i++ {
		stream.Send(&hello.HelloResponse{Message: fmt.Sprintf("%s %s %d", in.Name, "Reply", i)})
	}
	return nil
}

func RegisterToConsul(){
	consul.RegisterService("127.0.0.1:8500",&consul.ConsulService{
		Name: "helloworld",
		Tag:  []string{"helloworld"},
		IP:   "127.0.0.1",
		Port: 50051,
	})
}

type HealthImpl struct {
}
// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	fmt.Print("health checking\n")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}