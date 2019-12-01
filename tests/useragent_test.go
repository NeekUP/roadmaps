package tests

import (
	"github.com/NeekUP/roadmaps/core"
	"testing"
)

func TestUserAgentFingerprint(t *testing.T) {

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"MJ12bot", "Mozilla/5.0 (compatible; MJ12bot/v1.2.4; http://www.majestic12.co.uk/bot.php?+)", "  MJ12bot false"},
		{"GoogleBot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", "  Googlebot false"},
		{"Empty", "", "   false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := core.UserAgentFingerprint(tt.value); got != tt.want {
				t.Errorf("UserAgentFingerprint() = %v, want %v", got, tt.want)
			}
		})
	}
}
