# RocketMQ Go 工具包

这是一个 RocketMQ 的 Go 语言工具包，提供了简单易用的 API 来与 Apache RocketMQ 进行交互。

## 特性

- 支持同步、异步和单向消息发送
- 支持推模式消费
- 简洁的 API 设计
- 完善的错误处理
- 详细的示例代码

## 安装

```bash
go get github.com/qishenonly/rocketmq
```

## 快速开始

### 生产者示例

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/qishenonly/rocketmq"
)

func main() {
	// 创建配置
	config := rocketmq.NewDefaultConfig("127.0.0.1:9876", "producer_group")

	// 创建生产者
	p, err := rocketmq.NewProducer(config)
	if err != nil {
		log.Fatalf("创建生产者失败: %v", err)
	}

	// 启动生产者
	if err := p.Start(); err != nil {
		log.Fatalf("启动生产者失败: %v", err)
	}
	defer p.Shutdown()

	// 创建消息
	message := rocketmq.NewMessage("TestTopic", []byte("Hello RocketMQ!"))
	message.WithTag("TagA").WithKey("OrderID123")

	// 同步发送
	ctx := context.Background()
	msgID, err := p.SendSync(ctx, message)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}
	fmt.Printf("发送消息成功，消息ID: %s\n", msgID)
}
```

### 消费者示例

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/qishenonly/rocketmq"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	// 创建配置
	config := rocketmq.NewDefaultConfig("127.0.0.1:9876", "consumer_group")

	// 创建消费者
	c, err := rocketmq.NewPushConsumer(config)
	if err != nil {
		log.Fatalf("创建消费者失败: %v", err)
	}

	// 订阅主题
	err = c.Subscribe("TestTopic", "*", func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i, msg := range msgs {
			fmt.Printf("收到消息 %d: %s\n", i, string(msg.Body))
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		log.Fatalf("订阅主题失败: %v", err)
	}

	// 启动消费者
	if err := c.Start(); err != nil {
		log.Fatalf("启动消费者失败: %v", err)
	}

	// 等待信号退出
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	
	// 关闭消费者
	if err := c.Shutdown(); err != nil {
		log.Printf("关闭消费者失败: %v", err)
	}
}
```

## 配置选项

### 基本配置

```go
config := &rocketmq.Config{
    NameServer: "127.0.0.1:9876",  // RocketMQ 名称服务器地址
    GroupName:  "my_group",        // 生产者或消费者组名
    Timeout:    3 * time.Second,   // 超时时间
    RetryTimes: 3,                 // 重试次数
}
```

### 阿里云 RocketMQ 配置

```go
config := &rocketmq.Config{
    NameServer: "http://MQ_INST_XXXX.cn-hangzhou.mq.aliyuncs.com:80",
    GroupName:  "GID_XXXX",
    AccessKey:  "your_access_key",
    SecretKey:  "your_secret_key",
}
```

## 完整示例

完整的示例代码可以在 `examples` 目录下找到：

- 简单生产者: [examples/producer/main.go](examples/producer/main.go)
- 简单消费者: [examples/consumer/main.go](examples/consumer/main.go)
- 创建主题: [examples/create_topic/main.go](examples/create_topic/main.go)

## 注意事项

1. 确保 RocketMQ 服务器已经启动并且可以访问
2. 生产者和消费者的组名应该符合 RocketMQ 的命名规范
3. 在生产环境中，应该妥善处理所有错误情况
4. 在应用程序退出前，应该正确关闭生产者和消费者

## 贡献

欢迎提交 Issue 和 Pull Request 来帮助改进这个项目。

## 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。
