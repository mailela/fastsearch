package global

// Config 服务器设置
type Config struct {
	Addr          string `yaml:"addr"`          // 监听地址
	Data          string `json:"data"`          // 数据目录
	Negative_data string `json:"negative_data"` // 负面词数据目录
	Debug         bool   `yaml:"debug"`         // 调试模式
	AllowDrop     bool   `yaml:"allow_drop"`    // 允许删除数据库
	Dictionary    string `json:"dictionary"`    // 字典路径
	EnableAdmin   bool   `yaml:"enableAdmin"`   //启用admin
	Gomaxprocs    int    `json:"gomaxprocs"`    //GOMAXPROCS
	Shard         int    `yaml:"shard"`         //分片数
	Auth          string `json:"auth"`          //认证
	EnableGzip    bool   `yaml:"enableGzip"`    //是否开启gzip压缩
	Timeout       int64  `json:"timeout"`       //超时时间
	BufferNum     int    `yaml:"bufferNum"`     //分片缓冲数
}
