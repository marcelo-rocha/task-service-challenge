package msgpublisher

import (
	"errors"

	nats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var ErrNotConnected = errors.New("publisher is not connected")

type MsgPublisher struct {
	conn   *nats.Conn
	logger *zap.Logger
}

func New(logger *zap.Logger) *MsgPublisher {
	return &MsgPublisher{logger: logger}
}

func (p *MsgPublisher) Connect(url string) error {
	c, err := nats.Connect(url)
	if err != nil {
		return err
	}
	p.conn = c
	return nil
}

func (p *MsgPublisher) Publish(subject string, msg string) error {
	if p.conn == nil {
		return ErrNotConnected
	}
	return p.conn.Publish(subject, []byte(msg))
}

func (p *MsgPublisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
