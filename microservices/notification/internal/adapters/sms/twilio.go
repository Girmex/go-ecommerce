package sms

import (
	"context"
	"fmt"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioSMSSender struct {
	client    *twilio.RestClient
	fromPhone string
}

func NewTwilioSMSSender(accountSid, authToken, fromPhone string) ports.SMSSender {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &TwilioSMSSender{
		client:    client,
		fromPhone: fromPhone,
	}
}

func (s *TwilioSMSSender) Send(ctx context.Context, phone string, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(phone)
	params.SetFrom(s.fromPhone)
	params.SetBody(message)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS via Twilio: %w", err)
	}

	if resp.Sid == nil {
		return fmt.Errorf("twilio response missing SID")
	}

	return nil
}
