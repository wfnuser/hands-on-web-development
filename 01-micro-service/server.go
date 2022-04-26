package main

import (
	"context"
	"fmt"
	proto "go-micro-consul-demo/proto"

	rl "github.com/asim/go-micro/plugins/wrapper/ratelimiter/ratelimit/v4"
	consul "github.com/go-micro/plugins/v4/registry/consul"
	"github.com/juju/ratelimit"
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
	addrs := make([]string, 1)
	addrs[0] = "127.0.0.1:8500"

	r := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = addrs
	})

	limit := 5
	br := ratelimit.NewBucketWithRate(float64(limit), int64(limit))

	service := micro.NewService(
		micro.Name("my.micro.service"),
		micro.Registry(r),
		micro.WrapHandler(rl.NewHandlerWrapper(br, false)),
	)
	service.Init()
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))
	service.Run()
}
