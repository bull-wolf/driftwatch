package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/notify"
)

var sampleDrifts = []drift.Drift{
	{Service: "auth-service", Field: "image", Expected: "auth:v1.2", Actual: "auth:v1.0"},
	{Service: "auth-service", Field: "replicas", Expected: "3", Actual: "1"},
}

func TestSend_NoDrifts_DoesNothing(t *testing.T) {
	n := notify.New(notify.Config{Channel: notify.ChannelStdout})
	if err := n.Send(nil); err != nil {
		t.Fatalf("expected no error for empty drifts, got %v", err)
	}
}

func TestSend_StdoutChannel(t *testing.T) {
	n := notify.New(notify.Config{Channel: notify.ChannelStdout})
	if err := n.Send(sampleDrifts); err != nil {
		t.Fatalf("stdout send failed: %v", err)
	}
}

func TestSend_UnsupportedChannel(t *testing.T) {
	n := notify.New(notify.Config{Channel: "slack"})
	err := n.Send(sampleDrifts)
	if err == nil {
		t.Fatal("expected error for unsupported channel")
	}
}

func TestSend_WebhookChannel_Success(t *testing.T) {
	var received notify.Payload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := notify.New(notify.Config{
		Channel:    notify.ChannelWebhook,
		WebhookURL: server.URL,
	})
	if err := n.Send(sampleDrifts); err != nil {
		t.Fatalf("webhook send failed: %v", err)
	}
	if len(received.Drifts) != 2 {
		t.Errorf("expected 2 drifts in payload, got %d", len(received.Drifts))
	}
	if received.Timestamp == "" {
		t.Error("expected non-empty timestamp in payload")
	}
}

func TestSend_WebhookChannel_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := notify.New(notify.Config{
		Channel:    notify.ChannelWebhook,
		WebhookURL: server.URL,
	})
	err := n.Send(sampleDrifts)
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestSend_WebhookChannel_MissingURL(t *testing.T) {
	n := notify.New(notify.Config{Channel: notify.ChannelWebhook})
	err := n.Send(sampleDrifts)
	if err == nil {
		t.Fatal("expected error when webhook_url is empty")
	}
}
