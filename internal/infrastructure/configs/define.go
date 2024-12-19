package configs

var (
	Mode EnvMode // 开发环境
)

// EnvMode 开发环境
type EnvMode string

const (
	Development EnvMode = "dev" // 开发
	Production  EnvMode = "pro" // 生产
	Prerelease  EnvMode = "pre" // 预发布
)

type Bootstrap struct {
	Server   *Server `mapstructure:"server"`
	Log      *Log    `mapstructure:"log"`
	JWT      *JWT    `mapstructure:"jwt"`
	Data     *Data   `mapstructure:"data"`
	ConfPath *string `mapstructure:"conf_path"`
}
type Server struct {
	Port               int    `mapstructure:"port"`
	RateQPS            int    `mapstructure:"rate_qps"`
	TracerPort         int    `mapstructure:"tracer_port"`
	Name               string `mapstructure:"name"`
	MaxRequestBodySize int    `mapstructure:"max_request_body_size"`
}

type Log struct {
	OutPath    string `mapstructure:"out_path"`
	FilePrefix string `mapstructure:"file_prefix"`
	Level      int64  `mapstructure:"max_size"`
	MaxSize    int64  `mapstructure:"max_size"`
	MaxBackups int64  `mapstructure:"max_backups"`
	MaxAge     int64  `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type JWT struct {
	Issuer            string `mapstructure:"issuer"`
	SigningKey        string `mapstructure:"signing_key"`
	ExpirationToken   int64  `mapstructure:"expiration_token"`
	ExpirationRefresh int64  `mapstructure:"expiration_refresh"`
}
type Data struct {
	DataBase *DataBase `mapstructure:"database"`
	Redis    *Redis    `mapstructure:"redis"`
}

// DataBase 数据库
type DataBase struct {
	Driver       string `mapstructure:"driver"`
	Source       string `mapstructure:"source"`
	MaxIdleConns int32  `mapstructure:"max_idle_conns"`
	MaxOpenConns int32  `mapstructure:"max_open_conns"`
	LogLevel     int    `mapstructure:"log_level"`
}

// Redis 数据库
type Redis struct {
	Network      string `mapstructure:"network"`
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	Db           int64  `mapstructure:"db"`
	ReadTimeout  int64  `mapstructure:"read_timeout"`
	WriteTimeout int64  `mapstructure:"write_timeout"`
}
