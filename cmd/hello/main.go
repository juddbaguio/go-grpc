package main

import (
	"go-grpc/grpc/hello"
	"go-grpc/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)
	hello.RegisterHelloServiceServer(srv, &services.Hello{})

	log.Println("hello grpc-server is starting")
	if err := srv.Serve(lis); err != nil {
		srv.GracefulStop()
	}
}
