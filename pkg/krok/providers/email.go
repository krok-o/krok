package providers

// Event is an event type when a notification is sent out.
type Event string

// Payload is a payload that is sent by an email.
type Payload string

// Email defines an email sending provider's capabilities.
type Email interface {
	Notify(email string, event Event, payload Payload) error
}
