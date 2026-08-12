package sms

import (
	"context"
	"fmt"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

var _ ports.SMSSender = (*TwilioSMSSender)(nil)

type TwilioSMSSender struct {
	client *twilio.RestClient
	from   string
}

func NewTwilioSMSSender(
	accountSID string,
	authToken string,
	from string,
) *TwilioSMSSender {
	client := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: accountSID,
			Password: authToken,
		},
	)

	return &TwilioSMSSender{
		client: client,
		from:   from,
	}
}

func (s *TwilioSMSSender) Send(
	ctx context.Context,
	phone string,
	message string,
) error {
	params := &api.CreateMessageParams{}

	params.SetTo(phone)
	params.SetFrom(s.from)
	params.SetBody(message)

	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("send SMS with Twilio: %w", err)
	}

	return nil
}
