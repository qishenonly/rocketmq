package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/qishenonly/rocketmq"
)

func main() {
	// 创建配置
	config := rocketmq.NewDefaultConfig("127.0.0.1:9876", "simple_consumer_group")

	// 创建消费者
	c, err := rocketmq.NewPushConsumer(config)
	if err != nil {
		log.Fatalf("创建消费者失败: %v", err)
	}

	// 订阅主题
	err = c.Subscribe("TestTopic", "*", func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i, msg := range msgs {
			fmt.Printf("收到消息 %d: ID=%s, Tags=%s, Keys=%s, Body=%s\n",
				i, msg.MsgId, msg.GetTags(), msg.GetKeys(), string(msg.Body))
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

	fmt.Println("消费者已启动，按 Ctrl+C 退出")
	<-signalChan

	// 关闭消费者
	if err := c.Shutdown(); err != nil {
		log.Printf("关闭消费者失败: %v", err)
	}

	fmt.Println("消费者已关闭")
}
