package rocketmq

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConfig 表示配置无效
	ErrInvalidConfig = errors.New("invalid kafka configuration")

	// ErrProducerClosed 表示生产者已关闭
	ErrProducerClosed = errors.New("kafka producer is closed")

	// ErrConsumerClosed 表示消费者已关闭
	ErrConsumerClosed = errors.New("kafka consumer is closed")

	// ErrTopicEmpty 表示主题为空
	ErrTopicEmpty = errors.New("kafka topic cannot be empty")

	// ErrNoNameServer 表示没有配置名称服务器
	ErrNoNameServer = errors.New("no name server configured")

	// ErrNoGroupName 表示没有配置生产者或消费者组名
	ErrNoGroupName = errors.New("no group name configured")

	// ErrNoGroupID 表示没有配置消费者组ID
	ErrNoGroupID = errors.New("no consumer group ID configured")
)

// RocketMQError 是RocketMQ错误的包装
type RocketMQError struct {
	Op  string // 操作名称
	Err error  // 原始错误
}

// Error 实现error接口
func (e *RocketMQError) Error() string {
	if e.Op != "" {
		return fmt.Sprintf("rocketmq %s: %v", e.Op, e.Err)
	}
	return fmt.Sprintf("rocketmq error: %v", e.Err)
}

// Unwrap 返回原始错误
func (e *RocketMQError) Unwrap() error {
	return e.Err
}

// WrapError 包装错误
func WrapError(op string, err error) error {
	if err == nil {
		return nil
	}
	return &RocketMQError{Op: op, Err: err}
}
