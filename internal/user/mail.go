package user

import (
	"context"
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type IMailer interface {
	SendVerificationCode(ctx context.Context, emailTo, firstname, lastname, code string) error
}

type userMailer struct {
	client      *mailjet.Client
	appName     string
	senderEmail string
}

func NewMailer(client *mailjet.Client, appName, senderEmail string) IMailer {
	return &userMailer{
		client:      client,
		appName:     appName,
		senderEmail: senderEmail,
	}
}

func (m *userMailer) SendVerificationCode(ctx context.Context, emailTo, firstname, lastname, code string) error {
	msgs := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{Email: m.senderEmail, Name: m.appName},
			To: &mailjet.RecipientsV31{
				{Email: emailTo, Name: firstname + " " + lastname},
			},
			Subject:  m.appName + " verification code",
			TextPart: "Failed to load html.",
			HTMLPart: fmt.Sprintf("<h2>Welcome %s %s</h2><h3>Here is your verification code</h3><p>%s</p>", firstname, lastname, code),
		},
	}

	messages := mailjet.MessagesV31{Info: msgs}
	_, err := m.client.SendMailV31(&messages)
	if err != nil {
		return err
	}
	return nil
}
