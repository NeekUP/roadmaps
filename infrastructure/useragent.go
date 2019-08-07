package infrastructure

import (
	"fmt"
	"github.com/mssola/user_agent"
)

func UserAgentFingerprint(useragent string) string {
	ua := user_agent.New(useragent)
	brName, _ := ua.Browser()
	return fmt.Sprintf("%s %s %s %t", ua.Platform(), ua.OS(), brName, ua.Mobile())
}
