// Package notify implements drift notification dispatch for driftwatch.
//
// It supports multiple notification channels:
//
//   - stdout: prints drift summaries to standard output (useful for CI/CD pipelines)
//   - webhook: POSTs a JSON payload to a configured HTTP endpoint
//
// # Usage
//
//	cfg := notify.Config{
//		Channel:    notify.ChannelWebhook,
//		WebhookURL: "https://hooks.example.com/driftwatch",
//	}
//	n := notify.New(cfg)
//	if err := n.Send(drifts); err != nil {
//		log.Printf("notification failed: %v", err)
//	}
//
// When no drifts are present, Send is a no-op and returns nil.
// The webhook payload includes a UTC timestamp and the full list of drift records.
package notify
