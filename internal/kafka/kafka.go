package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// Producer Kafka 生产者
type Producer struct {
	writer *kafka.Writer
}

// NewProducer 创建 Kafka 生产者
func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

// SendMessage 发送消息
func (p *Producer) SendMessage(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
	}

	return p.writer.WriteMessages(ctx, msg)
}

// Close 关闭生产者
func (p *Producer) Close() error {
	return p.writer.Close()
}

// Consumer Kafka 消费者
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer 创建 Kafka 消费者
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

// ReadMessage 读取消息
func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// AuditLogMessage 审计日志消息
type AuditLogMessage struct {
	AdminID      int64  `json:"admin_id"`
	AdminUsername string `json:"admin_username"`
	Action       string `json:"action"`
	Resource     string `json:"resource"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	RequestBody  string `json:"request_body"`
	ResponseBody string `json:"response_body"`
	IP           string `json:"ip"`
	UserAgent    string `json:"user_agent"`
	Duration     int64  `json:"duration"`
	StatusCode   int    `json:"status_code"`
	CreatedAt    string `json:"created_at"`
}

// NewAuditLogProducer 创建审计日志生产者
func NewAuditLogProducer(brokers []string) *Producer {
	return NewProducer(brokers, "audit-logs")
}

// SendAuditLog 发送审计日志
func (p *Producer) SendAuditLog(ctx context.Context, auditLog *AuditLogMessage) error {
	log.Printf("Sending audit log to Kafka: %+v", auditLog)
	return p.SendMessage(ctx, auditLog.AdminUsername, auditLog)
}
