package mailgun

import (
	"bytes"
	"fmt"
	"time"

	mg "github.com/mailgun/mailgun-go"
	"github.com/rs/zerolog"

	"github.com/krok-o/krok/pkg/krok/providers"
)

var (
	// PasswordReset is an event that happens when the user's password is reset.
	PasswordReset providers.Event = "Password Reset"
	// GenerateConfirmCode is an event before password reset which sends a confirm code to the user's email address.
	GenerateConfirmCode providers.Event = "Confirm Code"
	// Welcome template for new sign-ups.
	Welcome providers.Event = "Welcome"
)

var (
	welcomeTemplate = `Dear %s
Thank you for signing up!`
	passwordResetTemplate = `Dear %s
Your password has been successfully reset to: %s. Please change as soon as possible.`
	confirmCodeTemplate = `Dear %s
Please enter the following code into the confirm code window: %s`
)

// Config is configuration for this provider.
type Config struct {
	Domain string
	APIKey string
}

// Dependencies contains dependencies this provider uses.
type Dependencies struct {
	Logger zerolog.Logger
}

// Sender is an email sender using Mailgun as a backend.
type Sender struct {
	Config
	Dependencies
}

// NewMailgunSender will connect to mailgun and return a mailgun email sender.
func NewMailgunSender(cfg Config, deps Dependencies) (*Sender, error) {
	return &Sender{
		Config:       cfg,
		Dependencies: deps,
	}, nil
}

// Notify attempts to send out an email using mailgun contaning the new password.
// Does not need to be a pointer receiver because it isn't storing anything.
func (e *Sender) Notify(email string, event providers.Event, payload providers.Payload) error {
	domain := e.Domain
	apiKey := e.APIKey
	sender := fmt.Sprintf("no-reply@%s", domain)
	subject := fmt.Sprintf("[%s] %s Notification", time.Now().Format("2006-01-02"), event)
	log := e.Logger.With().Str("email", email).Str("payload", string(payload)).Logger()

	if domain == "" && apiKey == "" {
		log.Warn().Msg("[WARNING] Mailgun not set up. Falling back to console output...")
		log.Info().Msg("A notification attempt was made for user.")
		return nil
	}

	var body string
	switch event {
	case PasswordReset:
		body = fmt.Sprintf(passwordResetTemplate, email, payload)
	case GenerateConfirmCode:
		body = fmt.Sprintf(confirmCodeTemplate, email, payload)
	case Welcome:
		body = fmt.Sprintf(welcomeTemplate, email)
	}

	mg := mg.NewMailgun(domain, apiKey)
	message := mg.NewMessage(sender, subject, body, email)
	_, _, err := mg.Send(message)
	return err
}

// BufferNotifier uses a byte buffer for notifications.
type BufferNotifier struct {
	buffer bytes.Buffer
}

// NewBufferNotifier creates a new notifier.
func NewBufferNotifier() *BufferNotifier {
	return &BufferNotifier{}
}

// Notify uses a buffer to store notifications for a user.
func (b *BufferNotifier) Notify(email string, event providers.Event, payload providers.Payload) error {
	var body string
	switch event {
	case PasswordReset:
		body = fmt.Sprintf(passwordResetTemplate, email, payload)
	case GenerateConfirmCode:
		body = fmt.Sprintf(confirmCodeTemplate, email, payload)
	}
	b.buffer.WriteString(body)
	return nil
}
