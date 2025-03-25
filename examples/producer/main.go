package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/qishenonly/rocketmq"
)

func main() {
	// 创建生产者
	p, err := rocketmq.NewProducer(
		rocketmq.WithNameServer("127.0.0.1:9876"),
		rocketmq.WithGroupName("simple_test_group"),
	)
	if err != nil {
		log.Fatalf("创建生产者失败: %v", err)
	}

	// 启动生产者
	if err := p.Start(); err != nil {
		log.Fatalf("启动生产者失败: %v", err)
	}
	defer p.Shutdown()

	// 创建消息
	message := rocketmq.NewMessage("test_topic_pc", []byte("Hello RocketMQ!"))
	message.WithTag("TagA").WithKey("OrderID123")

	// 同步发送
	ctx := context.Background()
	msgID, err := p.SendSync(ctx, message)
	if err != nil {
		log.Fatalf("同步发送消息失败: %v", err)
	}
	fmt.Printf("同步发送消息成功，消息ID: %s\n", msgID)

	// 异步发送
	for i := 0; i < 3; i++ {
		msg := rocketmq.NewMessage("test_topic_pc", []byte(fmt.Sprintf("Async Message %d", i)))
		msg.WithTag("TagB")

		err = p.SendAsync(ctx, msg, func(ctx context.Context, result *primitive.SendResult, err error) {
			if err != nil {
				fmt.Printf("异步发送消息失败: %v\n", err)
				return
			}
			fmt.Printf("异步发送消息成功，消息ID: %s\n", result.MsgID)
		})
		if err != nil {
			log.Printf("提交异步发送任务失败: %v", err)
		}
	}

	// 等待异步消息发送完成
	time.Sleep(time.Second)

	// 单向发送
	msg := rocketmq.NewMessage("test_topic_pc", []byte("One-way Message"))
	msg.WithTag("TagC")
	if err := p.SendOneWay(ctx, msg); err != nil {
		log.Printf("单向发送消息失败: %v", err)
	} else {
		fmt.Println("单向发送消息成功")
	}

	fmt.Println("所有消息发送完成")
}
