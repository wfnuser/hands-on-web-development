# go-micro-consul-demo

## 项目说明

本项目利用go-micro创建服务端(server.go)以及客户端(client.go)以演示微服务的基本架构。
本项目利用consul进行服务发现并以protobuf为通信协议。

在服务端中使用protobuf文件定义了一个服务叫做Greeter的处理器,它有一个接收HelloRequest并返回HelloResponse的Hello方法。
在客户端完成服务注册,客户端调用Hello,获取服务列表,获取单个服务,以及监听服务的功能。

## GRPC 接口定义

   ```protobuf
   syntax = "proto3";
   
   package helloworld;
   
   service Greeter {
       rpc Hello(HelloRequest) returns (HelloResponse) {}
   }
   
   message HelloRequest {
       string name = 1;
   }
   
   message HelloResponse {
       string greeting = 2;
   }
   ```

## server.go

   1. 服务端:使用go-micro创建服务端Demo,并注册到consul

      ```go
         registry := consul.NewRegistry(func(options *registry.Options) {
         		options.Addrs = addrs
         })
         service := micro.NewService(
         		// Set service name
         		micro.Name("my.micro.service"),
         		// Set service registry
         		micro.Registry(registry),
         )
         service.Run()
      
      
      ```

      

## client.go

   1. 客户端:使用go-micro创建客户端Demo,注册到consul.
      ```go
      	r := consul.NewRegistry(func(options *registry.Options) {
      		options.Addrs = addrs
      	})
      	service := micro.NewService(
      		micro.Name("my.micro.service.client"),
      		micro.Registry(r))
      ```

   2. 客户端rpc调用
      ```go
      	// 创建新的客户端
      	greeter := helloworld.NewGreeterService(serverName, service.Client())
      	// 调用greeter
      	rsp, err := greeter.Hello(context.TODO(), &helloworld.HelloRequest{Name: "John"})
      ```

   3. 查询服务列表
      ```go
      	services,err:=registry.ListServices()
      ```

   4. 获取某一个服务
      ```go
      	service, err := registry.GetService(serverName)
      ```

   5. 监听服务
      ```go
      	//监听服务
      	watch, err := registry.Watch()
      	for {
      		result, err := watch.Next()
      		if len(result.Action) > 0 {
      			fmt.Println(result, err)
      		}
      	}
      ```

## 如何运行

   0. 准备consul环境
      ```
      docker run -d -p 8500:8500 --restart=always --name=consul consul:latest agent -server -bootstrap -ui -node=1 -client='0.0.0.0'
      ```

   1. clone项目
      ```
      git clone git@github.com/wfnuser/web-development-in-action.git
      ```

   2. protoc
      ```
      go install go-micro.dev/v4/cmd/protoc-gen-micro@v4
      protoc  --micro_out=. --go_out=. proto/greeter.proto
      ```

   3. 启动服务端
      ```
      go run server.go --registry=consul --server_address=localhost:8500
      ```

   4. 启动客户端
      ```
      go run client.go
      ```

   5. 在consul中可以看到 客户端和服务端都已经注册在服务列表中

## 限流器
限流器通常也是微服务架构中的一个必不可少的环节；我们这里为限流器做一个简单的演示。

   1. 在客户端调整请求的频率
      ```
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
      ```

   2. 在服务端设置限流器允许的 rate
      ```
      limit := 5
      br := ratelimit.NewBucketWithRate(float64(limit), int64(limit))

      service := micro.NewService(
         micro.Name("my.micro.service"),
         micro.Registry(r),
         micro.WrapHandler(rl.NewHandlerWrapper(br, false)),
      )
      ```

   3. 分别尝试让客户端以每秒100次和每秒一次的频率请求服务端 观察现象