package controller

import (
	"context"
	"fmt"
	"go_gin_api/grpc/server/proto/hello"
	"go_gin_api/grpc/server/proto/listen"
)

type ListenController struct {
	
}
func (l *ListenController) ListenData(ctx context.Context, in *listen.Request) (*listen.Response, error) {
	return &listen.Response{Message : fmt.Sprintf("[%s]", in.Name)}, nil
}

func (l *ListenController) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{Message: in.Name},nil
}

func(l *ListenController)LotsOfReplies(in *hello.HelloRequest, stream hello.Hello_LotsOfRepliesServer) error {
	for i := 0; i < 10; i++ {
		stream.Send(&hello.HelloResponse{Message: fmt.Sprintf("%s %s %d", in.Name, "Reply", i)})
	}
	return nil
}