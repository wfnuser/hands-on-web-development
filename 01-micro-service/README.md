# go-micro-consul-demo

## 项目说明

此项目利用go-micro创建服务端(server.go)以及客户端(client.go).

1. 在服务端中使用protobuf文件定义了一个服务叫做Greeter的处理器,它有一个接收HelloRequest并返回HelloResponse的Hello方法。并将服务端注册到consul

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

2. 在客户端完成服务注册,客户端调用Hello,获取服务列表,获取单个服务,以及监听服务的功能。

## 功能说明

1. server.go

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

      

2. client.go

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

   ## Run

1. ```
   # clone项目
   git clone git@github.com/wfnuser/web-development-in-action.git
   ```

2. ```
   # 启动服务端
   go run server.go --registry=consul --server_address=localhost:8500
   
   ```

3. ```
   # 启动客户端
   go run client.go
   
   ```

4. 在 consul 中可以看到 客户端和服务端都已经注册在服务列表中。