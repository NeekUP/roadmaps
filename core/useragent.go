package core

import (
	"fmt"
	"github.com/mssola/user_agent"
)

func UserAgentFingerprint(useragent string) string {
	ua := user_agent.New(useragent)
	brName, _ := ua.Browser()
	ua.Bot()
	return fmt.Sprintf("platform:%s os:%s browser:%s mobile:%t bot:%v", ua.Platform(), ua.OS(), brName, ua.Mobile(), ua.Bot())
}
