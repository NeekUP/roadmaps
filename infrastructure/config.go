package infrastructure

type Config struct {
	Logger struct {
		Path string `json:"path"`
	}
	Db struct {
		ConnString string `json:"connString"`
	}
	//HTTPServer struct {
	//	StaticPath string `json:"staticpath"`
	//	Host       string `json:"host"`
	//	Schema     string `json:"schema"`
	//	Port       string `json:"port"`
	//}
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
