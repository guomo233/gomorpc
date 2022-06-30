# 设计

## RPC 协议

在 TCP 上层实现自定义 RPC 协议，不依赖 HTTP；

在一次连接中（在这 Header 是 RPC 协议自定义结构，与 HTTP 无关）：

```
|       Header        | MHeader1 | Message1 | MHeader2 | Message2 | ...
| <- 固定 JSON 编码 -> | <------ 编码方式由 Header.CodeType 决定 ----->
```

- Header：
  - MagicNumber：解码时校验读取的结构是 Header；
  - CodeType：消息编码方式（如 Protocol Buffer）；
- MHeader：
	- 用于请求：服务名、方法名、超时时间；
	- 用于响应：RPC 框架定义的异常，如 Timeout，函数的异常可以包含在返回值对象中；
	- 都有：请求ID（用于确定响应是哪个请求的）；
- Message：参数或返回值；

读取：

1. 反序列 MHeader，通过方法名配合反射确定参数类型；
2. 基于参数类型，进一步反序列化 Message；

写入：可以一次性将 MHeader 和 Message 写入（利用`bufio`一次性 flush，可以更高效）

连接包装：【编解码器【bufio【conn】】】

## 服务端

### 连接处理

采用 per goroutine per conn，其中每个 conn goroutine 又采用 per goroutine per request。由于不同的请求是基于一个连接的，所以发送响应需要考虑同步，否则响应消息会混杂

### 服务注册

`server`保存了`service name`到`service`的映射，`service`保存了`Method name`到`Method`的映射，同时还保存了对应的对象实例，因为在调用方法时需要将实例作为第一个参数传递进去

### 超时控制

对于调用请求处理函数，需要考虑超时控制。超时时间与客户端调用请求超时一致

### 简易注册中心

针对每个服务地址维护一个时间戳，记录上一次心跳的时间。当客户端需要返回服务地址时，检测其是否超时，超时则不返回该地址。避免了使用计时器

## 客户端

### 连接处理

同一连接的不同请求不能并发：

- 不能同时发送多个请求，否则内容会混再一起，只能串行发送；
- 不能同时读取响应，否则读出来的内容会乱，只能串行读，；

同一连接的不同请求异步：

- 存在一个 pending 列表（请求ID => done channel 的映射），发送请求的时候将请求存放至 pending，收到响应后将请求从 pending 中移除；
- `Call`只发送消息不等待响应，发送后获得一个 done channel，由后台 goroutine 读取响应，读取后根据请求ID找到 done channel 发送通知（发送 MHeader 返回的异常）；
- 当服务端或客户端发生异常时，将错误通知 pending 中的所有请求；

### 获取返回值

客户端必须以`Call(arg, ret)`的形式传递返回值，而不能`ret = Call(arg)`，因为`Call`函数不知道返回值的类型，所以无法自己创建出`ret`并返回，而以参数形式拿到的`ret`是用户调用`Call`之前自己创建好的

### 超时控制

需要控制超时的地方：

- 建立连接：可采用`net.DialTimeout`代替`net.Dial`；
- 调用请求：可由用户传入`context.WithTimeout`，处理`context.Done()`到期；

### 域名解析

采用和 gRPC 类似的方案：在客户端用一个全局的 resolver 解析域名为主机集合。需要在一开始将 resolver 注册进框架，可以注册多个，解析时按注册顺序解析

将 resolver 定义为接口，用户可自定义实现。框架提供了一个默认的静态解析实现，由用户自己传入域名与主机的映射

# Example

定义参数和返回值：

```protobuf
syntax = "proto3";
package proto;
option go_package="./";

message Arg {
	string msg = 1;
}

message Ret {
	string msg = 1;
}
```

运行服务器和客户端：

```go
package main

import (
	"fmt"
	"gomorpc"
	pb "main/proto"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// 定义 RPC 服务
type MyService struct {}

// 方法定义需按照如下规范
func (s *MyService) Func(arg *pb.Arg, ret *pb.Ret) {
	ret.Msg = arg.Msg + "--ret"
}

func main()  {
    // 运行服务器
	go func() {
		server := gomorpc.NewServer("0.0.0.0:63662")
		if err := server.Register(&MyService{}); err != nil {
			fmt.Println(err)
            return
		}
		
		if err := server.Server(); err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(time.Second)

    // 定义域名解析
	resolver := new(gomorpc.StaticResolver)
	resolver.Register(map[string][]string{
        // 域名需以 domain: 开头
		"domain:hostname": {"127.0.0.1:63662"},
	})
	gomorpc.RegisterResolver(resolver)
	
    // 创建客户端
    client, err := gomorpc.NewClient("domain:hostname", 
                                     gomorpc.PROTO_BUF, // 消息编码方式
                                     time.Second,       // 连接超时，0 表示无限制
                                     (gomorpc.RoundRabinLoadBalancer)) // 负载均衡
	if err != nil {
		fmt.Println(err)
		return
	}

    // 调用服务
    arg := &pb.Arg{Msg: s}
    ret := &pb.Ret{}
    done, err := client.Call("MyService", // 服务
                             "Func",      // 服务方法
                             arg,         // 参数
                             ret,         // 返回值
                             time.Second) // 调用超时，0 表示无限制
    if err != nil {
        fmt.Println(err)
        return
    }
	
    // 处理异步返回
    var wg sync.WaitGroup
    wg.Add(1)
    go func(ret *pb.Ret) {
        err := <- done // 返回 ERROR_TIMEOUT 表示超时，0 表示正常
        fmt.Println(err, ret.Msg)
        wg.Done()
    }(ret)
    wg.Wait()
}

```

