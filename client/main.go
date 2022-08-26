package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"go-grpc/grpc/hello"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewCredentials() credentials.TransportCredentials {
	b, err := ioutil.ReadFile("./infra/k8s/tls/tls.crt")
	if err != nil {
		log.Println(err)
		return nil
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		log.Println("failed to append")
		return nil
	}

	return credentials.NewTLS(&tls.Config{
		ServerName:         "juddbaguio.dev",
		RootCAs:            cp,
		InsecureSkipVerify: true,
	})
}

func main() {
	conn, err := grpc.Dial("juddbaguio.dev:443", grpc.WithTransportCredentials(NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := hello.NewHelloServiceClient(conn)

	res, err := client.SayHello(context.Background(), &hello.HelloRequest{
		Greeting: "Hello",
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println(res.Reply)

	repeatedRes, err := client.SayHelloRepeated(context.Background(), &hello.HelloRequestRepeated{
		Greeting: "WOW",
		Num:      23,
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {
		reply, err := repeatedRes.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(reply)
	}
}
