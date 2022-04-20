package main

import (
	"context"
	"fmt"
	proto "go-micro-nacos-demo/proto"

	consul "github.com/go-micro/plugins/v4/registry/consul"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	fmt.Println("Hello Received")
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	r := consul.NewRegistry(registry.Addrs("localhost:8500"))

	service := micro.NewService(
		micro.Name("my.micro.service"),
		micro.Registry(r),
	)
	service.Init()
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))
	service.Run()
}
