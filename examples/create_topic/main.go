package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qishenonly/rocketmq"
)

func main() {
	// 创建配置
	config := rocketmq.NewDefaultConfig("127.0.0.1:9876", "simple_producer_group")

	// 创建生产者
	p, err := rocketmq.NewProducer(config)
	if err != nil {
		log.Fatalf("创建生产者失败: %v", err)
	}

	// 创建主题配置
	topicConfig := &rocketmq.TopicConfig{
		Topic:          "NewTestTopic",
		ReadQueueNums:  8,
		WriteQueueNums: 8,
		Perm:           6, // 读写权限
	}

	// 创建主题
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 需要提供 Broker 地址
	brokerAddr := "127.0.0.1:10911" // 默认端口是 10911

	err = p.CreateTopic(ctx, brokerAddr, topicConfig)
	if err != nil {
		log.Fatalf("创建主题失败: %v", err)
	}

	fmt.Printf("成功创建主题: %s\n", topicConfig.Topic)
}
