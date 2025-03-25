package rocketmq

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// Producer RocketMQ 生产者接口
type Producer interface {
	// Start 启动生产者
	Start() error
	// Shutdown 关闭生产者
	Shutdown() error
	// SendSync 同步发送消息
	SendSync(ctx context.Context, message *Message) (string, error)
	// SendAsync 异步发送消息
	SendAsync(ctx context.Context, message *Message, callback func(context.Context, *primitive.SendResult, error)) error
	// SendOneWay 单向发送消息（不关心结果）
	SendOneWay(ctx context.Context, message *Message) error
	// CreateTopic 创建主题
	CreateTopic(ctx context.Context, brokerAddr string, config *TopicConfig) error
}

// DefaultProducer 默认生产者实现
type DefaultProducer struct {
	config   *Config
	producer rocketmq.Producer
	mutex    sync.Mutex
	started  bool
}

// NewProducer 创建新的生产者
func NewProducer(config *Config) (Producer, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &DefaultProducer{
		config:  config,
		started: false,
	}, nil
}

// Start 启动生产者
func (p *DefaultProducer) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.started {
		return nil
	}

	opts := []producer.Option{
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{p.config.NameServer})),
		producer.WithRetry(p.config.RetryTimes),
		producer.WithGroupName(p.config.GroupName),
	}

	// 如果提供了 AccessKey 和 SecretKey，则使用 ACL
	if p.config.AccessKey != "" && p.config.SecretKey != "" {
		opts = append(opts, producer.WithCredentials(primitive.Credentials{
			AccessKey: p.config.AccessKey,
			SecretKey: p.config.SecretKey,
		}))
	}

	var err error
	p.producer, err = rocketmq.NewProducer(opts...)
	if err != nil {
		return fmt.Errorf("create producer error: %w", err)
	}

	err = p.producer.Start()
	if err != nil {
		return fmt.Errorf("start producer error: %w", err)
	}

	p.started = true
	return nil
}

// Shutdown 关闭生产者
func (p *DefaultProducer) Shutdown() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.started {
		return nil
	}

	err := p.producer.Shutdown()
	if err != nil {
		return fmt.Errorf("shutdown producer error: %w", err)
	}

	p.started = false
	return nil
}

// convertMessage 将内部消息转换为 RocketMQ 消息
func convertMessage(message *Message) *primitive.Message {
	msg := primitive.NewMessage(message.Topic, message.Body)

	if message.Tags != "" {
		msg.WithTag(message.Tags)
	}

	if message.Keys != "" {
		msg.WithKeys([]string{message.Keys})
	}

	for k, v := range message.Properties {
		msg.WithProperty(k, v)
	}

	return msg
}

// SendSync 同步发送消息
func (p *DefaultProducer) SendSync(ctx context.Context, message *Message) (string, error) {
	if !p.started {
		if err := p.Start(); err != nil {
			return "", err
		}
	}

	msg := convertMessage(message)
	res, err := p.producer.SendSync(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("send sync message error: %w", err)
	}

	return res.MsgID, nil
}

// SendAsync 异步发送消息
func (p *DefaultProducer) SendAsync(ctx context.Context, message *Message, callback func(context.Context, *primitive.SendResult, error)) error {
	if !p.started {
		if err := p.Start(); err != nil {
			return err
		}
	}

	msg := convertMessage(message)
	err := p.producer.SendAsync(ctx, callback, msg)
	if err != nil {
		return fmt.Errorf("send async message error: %w", err)
	}

	return nil
}

// SendOneWay 单向发送消息
func (p *DefaultProducer) SendOneWay(ctx context.Context, message *Message) error {
	if !p.started {
		if err := p.Start(); err != nil {
			return err
		}
	}

	msg := convertMessage(message)
	err := p.producer.SendOneWay(ctx, msg)
	if err != nil {
		return fmt.Errorf("send one-way message error: %w", err)
	}

	return nil
}

// TopicConfig 主题配置
type TopicConfig struct {
	// Topic 主题名称
	Topic string
	// ReadQueueNums 读队列数量
	ReadQueueNums int
	// WriteQueueNums 写队列数量
	WriteQueueNums int
	// Perm 权限，默认为读写权限
	Perm int
}

// CreateTopic 创建主题
func (p *DefaultProducer) CreateTopic(ctx context.Context, brokerAddr string, config *TopicConfig) error {
	// 创建管理客户端
	adminClient, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver([]string{p.config.NameServer})))
	if err != nil {
		return fmt.Errorf("create admin client error: %w", err)
	}
	defer adminClient.Close()

	// 创建主题
	err = adminClient.CreateTopic(
		ctx,
		admin.WithTopicCreate(config.Topic),
		admin.WithBrokerAddrCreate(brokerAddr),
		admin.WithReadQueueNums(config.ReadQueueNums),
		admin.WithWriteQueueNums(config.WriteQueueNums),
		admin.WithPerm(config.Perm),
	)
	if err != nil {
		return fmt.Errorf("create topic error: %w", err)
	}

	return nil
}
