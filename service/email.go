package service

import (
	"context"

	"github.com/wneessen/go-mail"
)

type Email struct {
	client       *mail.Client
	sender       string
	company      string
	web_hostname string
}

func NewEmailService(client *mail.Client, sender string, company string, web_hostname string) Email {
	return Email{
		client:       client,
		sender:       sender,
		company:      company,
		web_hostname: web_hostname,
	}
}

func (es Email) DialAndSend(ctx context.Context, ml ...*mail.Msg) error {
	return es.client.DialAndSendWithContext(ctx, ml...)
}

func (es Email) GetSenderEmail() string {
	return es.sender
}

func (es Email) GetCompanyEmail() string {
	return es.company
}

func (es Email) GetWebHostname() string {
	return es.web_hostname
}
