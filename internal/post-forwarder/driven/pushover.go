package driven

import (
	"context"
	"fmt"

	"github.com/adlandh/post-forwarder/internal/post-forwarder/domain"
	"github.com/gregdel/pushover"
)

var _ domain.Notifier = (*Pushover)(nil)

type Pushover struct {
	client     *pushover.Pushover
	recipients []pushover.Recipient
}

// NewPushover returns a new instance of a Pushover notification service.
// For more information about Pushover app token:
//
//	-> https://support.pushover.net/i175-how-do-i-get-an-api-or-application-token
func NewPushover(appToken string) *Pushover {
	return &Pushover{
		client:     pushover.New(appToken),
		recipients: []pushover.Recipient{},
	}
}

// AddReceivers takes Pushover user/group IDs and adds them to the internal recipient list. The Send method will send
// a given message to all of those recipients.
func (p *Pushover) AddReceivers(recipientIDs ...string) {
	for _, recipient := range recipientIDs {
		p.recipients = append(p.recipients, *pushover.NewRecipient(recipient))
	}
}

// Send takes a message subject and a message body and sends them to all previously set recipients.
func (p *Pushover) Send(ctx context.Context, subject, message string) error {
	for i := range p.recipients {
		select {
		case <-ctx.Done():
			return fmt.Errorf("error sending message to Pushover recipient '%s': %w", p.recipients[i], ctx.Err())
		default:
			_, err := p.client.SendMessage(
				&pushover.Message{
					Message: message,
					Title:   subject,
					HTML:    true,
				},
				&p.recipients[i],
			)
			if err != nil {
				return fmt.Errorf("failed to send message to Pushover recipient '%s': %w", p.recipients[i], err)
			}
		}
	}

	return nil
}
