package rocketmq

import (
	"fmt"
	"time"
)

// Config 包含 RocketMQ 的基本配置
type Config struct {
	// NameServer RocketMQ 名称服务器地址，格式为 "ip:port"，多个地址用分号分隔
	NameServer string
	// GroupName 生产者或消费者组名
	GroupName string
	// AccessKey 阿里云 RocketMQ 的 AccessKey（如果使用阿里云 RocketMQ）
	AccessKey string
	// SecretKey 阿里云 RocketMQ 的 SecretKey（如果使用阿里云 RocketMQ）
	SecretKey string
	// Timeout 超时时间
	Timeout time.Duration
	// RetryTimes 重试次数
	RetryTimes int
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig(nameServer, groupName string) *Config {
	return &Config{
		NameServer: nameServer,
		GroupName:  groupName,
		Timeout:    3 * time.Second,
		RetryTimes: 3,
	}
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	if c.NameServer == "" {
		return fmt.Errorf("name server address cannot be empty")
	}
	if c.GroupName == "" {
		return fmt.Errorf("group name cannot be empty")
	}
	return nil
}

// Message 消息结构体
type Message struct {
	// Topic 消息主题
	Topic string
	// Tags 消息标签
	Tags string
	// Keys 消息键
	Keys string
	// Body 消息体
	Body []byte
	// Properties 消息属性
	Properties map[string]string
}

// NewMessage 创建新消息
func NewMessage(topic string, body []byte) *Message {
	return &Message{
		Topic:      topic,
		Body:       body,
		Properties: make(map[string]string),
	}
}

// WithTag 设置消息标签
func (m *Message) WithTag(tag string) *Message {
	m.Tags = tag
	return m
}

// WithKey 设置消息键
func (m *Message) WithKey(key string) *Message {
	m.Keys = key
	return m
}

// WithProperty 设置消息属性
func (m *Message) WithProperty(key, value string) *Message {
	m.Properties[key] = value
	return m
}

// NewDefaultTopicConfig 创建默认主题配置
func NewDefaultTopicConfig(topic string) *TopicConfig {
	return &TopicConfig{
		Topic:          topic,
		ReadQueueNums:  4,
		WriteQueueNums: 4,
		Perm:           6, // 读写权限
	}
}
