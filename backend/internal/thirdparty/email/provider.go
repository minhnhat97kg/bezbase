package email

import "context"

// EmailProvider defines the interface for email sending providers
type EmailProvider interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	IsConfigured() bool
}