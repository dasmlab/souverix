package sip

import (
	"strings"
	"testing"
)

func TestMessage_IsRequest(t *testing.T) {
	tests := []struct {
		name    string
		message *Message
		want    bool
	}{
		{
			name: "INVITE request",
			message: &Message{
				Method:  MethodINVITE,
				URI:     "sip:alice@example.com",
				Version: "SIP/2.0",
			},
			want: true,
		},
		{
			name: "200 OK response",
			message: &Message{
				Version:    "SIP/2.0",
				StatusCode: StatusOK,
				StatusText: "OK",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.message.IsRequest(); got != tt.want {
				t.Errorf("Message.IsRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_GetHeader(t *testing.T) {
	msg := &Message{
		Headers: map[string][]string{
			"From": {"sip:alice@example.com"},
			"Via":  {"SIP/2.0/UDP 192.168.1.1:5060"},
		},
	}

	tests := []struct {
		name     string
		header   string
		want     string
		wantCase bool
	}{
		{"exact case", "From", "sip:alice@example.com", true},
		{"lowercase", "from", "sip:alice@example.com", true},
		{"uppercase", "FROM", "sip:alice@example.com", true},
		{"mixed case", "ViA", "SIP/2.0/UDP 192.168.1.1:5060", true},
		{"missing", "To", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := msg.GetHeader(tt.header)
			if (got != "") != tt.wantCase {
				t.Errorf("Message.GetHeader() = %v, wantCase %v", got, tt.wantCase)
			}
			if tt.wantCase && got != tt.want {
				t.Errorf("Message.GetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_SetHeader(t *testing.T) {
	msg := &Message{
		Headers: make(map[string][]string),
	}

	msg.SetHeader("From", "sip:alice@example.com")

	if got := msg.GetHeader("From"); got != "sip:alice@example.com" {
		t.Errorf("Message.SetHeader() failed, got %v", got)
	}
}

func TestMessage_String(t *testing.T) {
	msg := &Message{
		Method:  MethodINVITE,
		URI:     "sip:bob@example.com",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"From": {"sip:alice@example.com"},
			"To":   {"sip:bob@example.com"},
		},
		Body: "v=0\r\no=alice 2890844526 2890844526 IN IP4 192.168.1.1",
	}

	result := msg.String()

	if !strings.Contains(result, "INVITE") {
		t.Error("Message.String() missing method")
	}
	if !strings.Contains(result, "From:") {
		t.Error("Message.String() missing From header")
	}
	if !strings.Contains(result, "v=0") {
		t.Error("Message.String() missing body")
	}
}

func TestMessage_AddHeader(t *testing.T) {
	msg := &Message{
		Headers: make(map[string][]string),
	}

	msg.AddHeader("Via", "SIP/2.0/UDP 192.168.1.1:5060")
	msg.AddHeader("Via", "SIP/2.0/UDP 192.168.1.2:5060")

	all := msg.GetHeaderAll("Via")
	if len(all) != 2 {
		t.Errorf("Message.AddHeader() failed, got %d values, want 2", len(all))
	}
}
