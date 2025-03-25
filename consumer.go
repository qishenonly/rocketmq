package rocketmq

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

// MessageHandler 消息处理函数类型
type MessageHandler func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)

// Consumer RocketMQ 消费者接口
type Consumer interface {
	// Start 启动消费者
	Start() error
	// Shutdown 关闭消费者
	Shutdown() error
	// Subscribe 订阅主题
	Subscribe(topic string, expression string, handler MessageHandler) error
}

// PushConsumer 推模式消费者
type PushConsumer struct {
	config   *Config
	consumer rocketmq.PushConsumer
	mutex    sync.Mutex
	started  bool
}

// NewPushConsumer 创建推模式消费者
func NewPushConsumer(opts ...ConfigOption) (Consumer, error) {
	// 使用默认配置
	config := DefaultConfig()

	// 应用选项
	for _, opt := range opts {
		opt(config)
	}

	return &PushConsumer{
		config:  config,
		started: false,
	}, nil
}

// Start 启动消费者
func (c *PushConsumer) Start() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.started {
		return nil
	}

	opts := []consumer.Option{
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{c.config.NameServer})),
		consumer.WithGroupName(c.config.GroupName),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
		consumer.WithAutoCommit(true),
	}

	// 如果提供了 AccessKey 和 SecretKey，则使用 ACL
	if c.config.AccessKey != "" && c.config.SecretKey != "" {
		opts = append(opts, consumer.WithCredentials(primitive.Credentials{
			AccessKey: c.config.AccessKey,
			SecretKey: c.config.SecretKey,
		}))
	}

	var err error
	c.consumer, err = rocketmq.NewPushConsumer(opts...)
	if err != nil {
		return fmt.Errorf("create consumer error: %w", err)
	}

	err = c.consumer.Start()
	if err != nil {
		return fmt.Errorf("start consumer error: %w", err)
	}

	c.started = true
	return nil
}

// Shutdown 关闭消费者
func (c *PushConsumer) Shutdown() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.started {
		return nil
	}

	err := c.consumer.Shutdown()
	if err != nil {
		return fmt.Errorf("shutdown consumer error: %w", err)
	}

	c.started = false
	return nil
}

// Subscribe 订阅主题
func (c *PushConsumer) Subscribe(topic string, expression string, handler MessageHandler) error {
	if !c.started {
		if err := c.Start(); err != nil {
			return err
		}
	}

	err := c.consumer.Subscribe(topic, consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: expression,
	}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		return handler(ctx, msgs...)
	})

	if err != nil {
		return fmt.Errorf("subscribe topic error: %w", err)
	}

	return nil
}

// PullConsumer 拉模式消费者
type PullConsumer struct {
	config   *Config
	consumer rocketmq.PullConsumer
	mutex    sync.Mutex
	started  bool
}

// NewPullConsumer 创建拉模式消费者
func NewPullConsumer(config *Config) (Consumer, error) {
	return &PullConsumer{
		config:  config,
		started: false,
	}, nil
}

// Start 启动消费者
func (c *PullConsumer) Start() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.started {
		return nil
	}

	// 目前 RocketMQ Go 客户端不直接支持 PullConsumer
	// 这里是一个简化的实现
	return fmt.Errorf("pull consumer is not fully implemented in current version")
}

// Shutdown 关闭消费者
func (c *PullConsumer) Shutdown() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.started {
		return nil
	}

	c.started = false
	return nil
}

// Subscribe 订阅主题
func (c *PullConsumer) Subscribe(topic string, expression string, handler MessageHandler) error {
	return fmt.Errorf("pull consumer does not support subscribe method")
}
