package models

import (
	"encoding/json"
	"testing"
)

func TestConfidenceLevelString(t *testing.T) {
	tests := []struct {
		name  string
		level ConfidenceLevel
		want  string
	}{
		{name: "low", level: ConfidenceLow, want: "low"},
		{name: "medium", level: ConfidenceMedium, want: "medium"},
		{name: "high", level: ConfidenceHigh, want: "high"},
		{name: "unknown", level: ConfidenceLevel(99), want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConfidenceLevelMarshalJSON(t *testing.T) {
	b, err := json.Marshal(ConfidenceMedium)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(b) != `"medium"` {
		t.Fatalf("Marshal = %s, want %q", string(b), `"medium"`)
	}
}

func TestConfidenceLevelUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    ConfidenceLevel
		wantErr bool
	}{
		{name: "low", payload: `"low"`, want: ConfidenceLow},
		{name: "medium", payload: `"medium"`, want: ConfidenceMedium},
		{name: "high", payload: `"high"`, want: ConfidenceHigh},
		{name: "unknown defaults low", payload: `"other"`, want: ConfidenceLow},
		{name: "invalid type", payload: `10`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ConfidenceLevel
			err := json.Unmarshal([]byte(tt.payload), &got)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if got != tt.want {
				t.Fatalf("Unmarshal got %v, want %v", got, tt.want)
			}
		})
	}
}
