package misc

import (
	"os"
	"time"

	logrus "github.com/sirupsen/logrus"
	//"go.elastic.co/apm"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	RefreshInterval time.Duration
	HTTP            HTTPConfig
	Log             LogConfig
	Redis         RedisConfig
	RedisCluster  RedisClusterConfig
}

type HTTPConfig struct {
	Address        string
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	Domain         string
	AuctioneerHost string
	BidderHost     string
	CookieExpiry   int
}

// LogConfig Logging configuration
type LogConfig struct {
	Level      string
	Format     string
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

// RedisClusterConfig Redis configuration parameters
type RedisClusterConfig struct {
	Master     string
	Replica    string
	Password   string
	DB         int
	MaxRetries int
	Expiration time.Duration
}

// RedisConfig Redis configuration parameters
type RedisConfig struct {
	Address    string
	Password   string
	DB         int
	MaxRetries int
	Expiration time.Duration
}

// InitLogging Initialize logging framework
func InitLogging(lc *LogConfig) {

	var log = &logrus.Logger{}

	switch lc.Format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "log.level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function.name", // non-ECS
			},
		})
	default:
		fallthrough
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	switch lc.Level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	default:
		fallthrough
	case "info":
		log.SetLevel(logrus.InfoLevel)
	}

	if lc.Filename == "" {
		log.SetOutput(os.Stdout)
	} else {
		l := &lumberjack.Logger{
			Filename:   lc.Filename,
			MaxSize:    lc.MaxSize,
			MaxAge:     lc.MaxAge,
			MaxBackups: lc.MaxBackups,
			LocalTime:  lc.LocalTime,
			Compress:   lc.Compress,
		}

		log.SetOutput(l)
	}
}
