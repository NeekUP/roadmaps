package infrastructure

type Config struct {
	SiteHost string `json:"siteHost"`
	Db       DbConf
	Logger   struct {
		Path string `json:"path"`
	}
	ImgSaver struct {
		LocalFolder string `json:"localFolder"`
		UriPath     string `json:"uriPath"`
	}
	HTTPServer struct {
		StaticPath            string `json:"staticpath"`
		Host                  string `json:"host"`
		Schema                string `json:"schema"`
		Port                  string `json:"port"`
		ReadTimeoutSec        int    `json:"readTimeoutSec"`
		WriteTimeoutSec       int    `json:"writeTimeoutSec"`
		ReadHeadersTimeoutSec int    `json:"readHeadersTimeoutSec"`
	}
	Client struct {
		Host string `json:"host"`
	}
	SMTP struct {
		SenderEmail string `json:"senderEmail"`
		SenderName  string `json:"senderName"`
		Host        string `json:"host"`
		Port        int    `json:"port"`
		Pass        string `json:"pass"`
	}
	OAuth struct {
		ReturnUrl string           `json:"returnUrl"`
		Providers []OauthProviders `json:"providers"`
	}
	//Cache struct {
	//	Enable    bool   `json:"enable"`
	//	host      string `json:"host"`
	//	port      string `json:"port"`
	//	Password  string `json:"password"`
	//	User      string `json:"user"`
	//	DB        string `json:"db"`
	//	PoolSize  int    `json:"poolsize"`
	//	Duraction struct {
	//		Categories int `json:"categories"`
	//		Category   int `json:"category"`
	//		Feed       int `json:"feed"`
	//		Channels   int `json:"channels"`
	//		Post       int `json:"post"`
	//		Tag        int `json:"tag"`
	//	}
	//}
}

type OauthProviders struct {
	Name     string   `json:"name"`
	ClientId string   `json:"clientId"`
	Secret   string   `json:"secret"`
	Scope    []string `json:"scope"`
}

type DbConf struct {
	ConnString string `json:"connString"`
	// Valid levels:
	//	trace
	//	debug
	//	info
	//	warn
	//	error
	//	none
	LogLevel string `json:"logLevel"`
}
