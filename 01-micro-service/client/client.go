package main

import (
	"context"
	"fmt"
	proto "go-micro-consul-demo/proto"
	"time"

	consul "github.com/go-micro/plugins/v4/registry/consul"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

const serverName = "my.micro.service"

func main() {
	addrs := make([]string, 1)
	addrs[0] = "localhost:8500"

	registry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = addrs
	})

	// 定义服务，可以传入其它可选参数
	service := micro.NewService(
		micro.Name("my.micro.service.client"),
		micro.Registry(registry))

	// 获取所有服务
	fmt.Println(registry.ListServices())

	// 创建新的客户端
	greeter := proto.NewGreeterService(serverName, service.Client())
	// 调用greeter

	ticker := time.NewTicker(1 * time.Second / 100)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "John"})
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(rsp)
		}
	}()

	// 监听服务
	watch, err := registry.Watch()
	if err != nil {
		fmt.Println(err)
	}

	// 打印响应请求
	// fmt.Println(rsp.Greeting)
	go service.Run()

	for {
		result, err := watch.Next()
		if len(result.Action) > 0 {
			fmt.Println(result, err)
		}
	}
}
