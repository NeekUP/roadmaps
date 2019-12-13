package infrastructure

type Config struct {
	Db     DbConf
	Logger struct {
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
	//Cache struct {
	//	Enable    bool   `json:"enable"`
	//	Host      string `json:"host"`
	//	Port      string `json:"port"`
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
