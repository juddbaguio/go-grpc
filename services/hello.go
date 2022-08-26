package services

import (
	"context"
	"fmt"
	"go-grpc/grpc/hello"
	"log"
	"time"
)

type Hello struct {
	hello.UnimplementedHelloServiceServer
}

func (h *Hello) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	log.Println(req.Greeting)
	return &hello.HelloReply{
		Reply: "You're welcome!",
	}, nil
}
func (h *Hello) SayHelloRepeated(req *hello.HelloRequestRepeated, stream hello.HelloService_SayHelloRepeatedServer) error {
	log.Println(req.Greeting, req.Num)

	for i := 0; i < int(req.Num); i++ {
		err := stream.Send(&hello.HelloReply{
			Reply: fmt.Sprintf("%v - %v", req.Greeting, i),
		})

		if err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
