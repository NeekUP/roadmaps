package tests

import (
	"roadmaps/core"
	"testing"
)

func TestUserAgentFingerprint(t *testing.T) {

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := core.UserAgentFingerprint(tt.value); got != tt.want {
				t.Errorf("UserAgentFingerprint() = %v, want %v", got, tt.want)
			}
		})
	}
}
