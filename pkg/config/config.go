package cfg

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"time"
)

// ServiceConfiguration  configuration for service
type ServiceConfiguration struct {
	Port string
}

// MySQLConfiguration  configuration for MySQL database connection
type MySQLConfiguration struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	LogMode  MySQLLogMode
	Charset  string
}

// PostgresConfiguration  configuration for Postgres database connection
type PostgresConfiguration struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	LogMode  MySQLLogMode
}

// ESConfiguration  configuration for elasticsearch connection
type ESConfiguration struct {
	Host                         string
	User                         string
	Password                     string
	ResponseHeaderTimeoutSeconds int
}

// RedisConfiguration ...
type RedisConfiguration struct {
	Host string
	Port string
}

// MongoDBConfiguration  configuration for redis connection
type MongoDBConfiguration struct {
	Host   string
	DBName string
}

type InfluxConfiguration struct {
	Host  string
	Port  string
	Token string
}

// KafkaConfiguration ...
type KafkaConfiguration struct {
	Brokers    []string
	User       string
	Password   string
	EnableSASL bool
}

// MinIOConfiguration ...
type MinIOConfiguration struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Token     string
	SSL       bool
	Region    string
}

// CreeperConfiguration ...
type CreeperConfiguration struct {
	Uri   string
	Token string
}

// WarehouseConfiguration ...
type WarehouseConfiguration struct {
	Uri       string
	AccessKey string
	SecretKey string
	Timeout   time.Duration
}

// LoggerConfig configuration for logger
type LoggerConfig struct {
	// log filename, **if it's not set, the log will be written to os.Stdout**
	Filename string
	// maximun size of single file, unit: MB
	MaxSize int
	// maximun number of days to retain the log file, unit: day
	MaxAge int
	// maximum number of log files to retain
	MaxBackups int
	// whether compress log file
	Compress bool
	// ??????????????? "warn" > "info" > "debug" > "trace"
	Level LevelMode
}

// RemoteServiceConfiguration  ...
type RemoteServiceConfiguration struct {
	// eg: "180.96.8.140:9998"
	RemoteHost string
	// eg: "/gateway/twirp/rime.adam.rpc_gateway.RPCGateway/"
	RemotePathPrefix string
}

// InitConfiguration ...
func InitConfiguration(configName string, configPaths []string, config interface{}) error {
	vp := viper.New()
	vp.SetConfigName(configName)
	vp.AutomaticEnv()
	for _, configPath := range configPaths {
		vp.AddConfigPath(configPath)
	}

	if err := vp.ReadInConfig(); err != nil {
		return errors.WithStack(err)
	}

	err := vp.Unmarshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// LevelMode ...
type LevelMode string

// LevelMode ...
// warn > info > debug > trace
const (
	LevelWarn  LevelMode = "warn"
	LevelInfo  LevelMode = "info"
	LevelDebug LevelMode = "debug"
	LevelTrace LevelMode = "trace"
)

// Level ...
func (l LevelMode) Level() logrus.Level {
	switch l {
	case LevelWarn:
		return logrus.WarnLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelDebug:
		return logrus.DebugLevel
	case LevelTrace:
		return logrus.TraceLevel
	}
	return logrus.WarnLevel
}

// IsDebugMode ...
func (l LevelMode) IsDebugMode() bool {
	return l.Level() >= logrus.DebugLevel
}

// MySQLLogMode ...
type MySQLLogMode string

// Console ?????? gorm ??? logger??????????????????sql????????????
// SlowQuery ??????????????? logger.Logger,???????????????sql?????????
// None ?????? log ??????
const (
	Console   MySQLLogMode = "console"
	SlowQuery MySQLLogMode = "slow_query"
	None      MySQLLogMode = "none"
)
