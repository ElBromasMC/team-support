package service

import (
	"context"

	"github.com/wneessen/go-mail"
)

type Email struct {
	client *mail.Client
}

func NewEmailService(client *mail.Client) Email {
	return Email{
		client: client,
	}
}

func (es Email) DialAndSend(ctx context.Context, ml ...*mail.Msg) error {
	return es.client.DialAndSendWithContext(ctx, ml...)
}
