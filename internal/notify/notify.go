// Package notify provides notification support for drift events,
// allowing alerts to be sent to configured channels such as stdout,
// webhook endpoints, or log sinks.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Channel represents a notification destination.
type Channel string

const (
	ChannelStdout  Channel = "stdout"
	ChannelWebhook Channel = "webhook"
)

// Config holds configuration for the notifier.
type Config struct {
	Channel    Channel `yaml:"channel"`
	WebhookURL string  `yaml:"webhook_url,omitempty"`
}

// Payload is the structure sent to webhook endpoints.
type Payload struct {
	Timestamp string       `json:"timestamp"`
	Drifts    []drift.Drift `json:"drifts"`
}

// Notifier sends drift notifications to a configured channel.
type Notifier struct {
	cfg    Config
	client *http.Client
}

// New creates a new Notifier with the given configuration.
func New(cfg Config) *Notifier {
	return &Notifier{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Send dispatches drift results to the configured channel.
// Returns an error if the channel is unsupported or the send fails.
func (n *Notifier) Send(drifts []drift.Drift) error {
	if len(drifts) == 0 {
		return nil
	}
	switch n.cfg.Channel {
	case ChannelStdout:
		return n.sendStdout(drifts)
	case ChannelWebhook:
		return n.sendWebhook(drifts)
	default:
		return fmt.Errorf("notify: unsupported channel %q", n.cfg.Channel)
	}
}

func (n *Notifier) sendStdout(drifts []drift.Drift) error {
	for _, d := range drifts {
		fmt.Printf("[DRIFT] service=%s field=%s want=%s got=%s\n",
			d.Service, d.Field, d.Expected, d.Actual)
	}
	return nil
}

func (n *Notifier) sendWebhook(drifts []drift.Drift) error {
	if n.cfg.WebhookURL == "" {
		return fmt.Errorf("notify: webhook_url is required for webhook channel")
	}
	payload := Payload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Drifts:    drifts,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}
	resp, err := n.client.Post(n.cfg.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: webhook post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
