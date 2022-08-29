package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"go-grpc/grpc/auth"
	"go-grpc/grpc/hello"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewCredentials(serverName string) credentials.TransportCredentials {
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
		ServerName:         serverName,
		RootCAs:            cp,
		InsecureSkipVerify: true,
	})
}

func main() {
	conn, err := grpc.Dial("hello.juddbaguio.dev:443", grpc.WithTransportCredentials(NewCredentials("hello.juddbaguio.dev")))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	authConn, err := grpc.Dial("auth.juddbaguio.dev:443", grpc.WithTransportCredentials(NewCredentials("auth.juddbaguio.dev")))
	if err != nil {
		log.Fatal(err)
	}
	defer authConn.Close()

	log.Println("connected successfully")
	helloSrv := hello.NewHelloServiceClient(conn)
	res, err := helloSrv.SayHello(context.Background(), &hello.HelloRequest{
		Greeting: "Hello",
	})

	if err != nil {
		log.Println("GRPC ERROR: ", err.Error())
		os.Exit(1)
	}

	log.Println(res.Reply)

	repeatedRes, err := helloSrv.SayHelloRepeated(context.Background(), &hello.HelloRequestRepeated{
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
	authSrv := auth.NewAuthServiceClient(authConn)

	loginRes, err := authSrv.HandleLogin(context.Background(), &auth.Login{
		Username: "Hello!",
		Password: "WOW",
	})

	if err != nil {
		log.Println("error: ", err.Error())
		os.Exit(1)
	}

	log.Println(loginRes)
}
