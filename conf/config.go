package conf

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/go-ini/ini"
)

const (
	appINIFilePath = "conf/app.ini"
)

var (
	Server     *Service
	DBConf     *Database
	RedisConf  *RedisConfig
	LogConf    *LogConfig
	ConsulConf *ConsulConfig
	EurekaConf *EurekaConfig
)

type LogConfig struct {
	LogLevel   string
	FileName   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type Service struct {
	Port   int
	Source string
}

type Database struct {
	DriverName string
	User       string
	Password   string
	DBHostname string
	DBPort     string
	DBName     string
}

type ConsulConfig struct {
	ServiceHost   string
	ConsulAddress string
	ServiceName   string
	CheckTimeout  string
	CheckInterval string
}

type EurekaConfig struct {
	EurekaURL   string
	LocalIP     string
	ServiceName string
}

func LoadConfig(confFilePath string) {
	if confFilePath == "" {
		confFilePath = appINIFilePath
	}
	absPath, err := filepath.Abs(confFilePath)
	if err != nil {
		panic(err)
	}
	cfg, err := ini.Load(absPath)
	if err != nil {
		panic(err)
	}
	LogConf = new(LogConfig)
	DBConf = new(Database)
	RedisConf = new(RedisConfig)
	Server = new(Service)
	ConsulConf = new(ConsulConfig)
	EurekaConf = new(EurekaConfig)
	mapTo("Log", LogConf, cfg)
	mapTo("Database", DBConf, cfg)
	mapTo("Redis", RedisConf, cfg)
	mapTo("Server", Server, cfg)
	mapTo("ConsulConfig", ConsulConf, cfg)
	mapTo("EurekaConfig", EurekaConf, cfg)

	fmt.Println("[Edgex-gateway] Config load finished!")
}

func mapTo(section string, v interface{}, cfg *ini.File) {
	if cfg == nil || section == "" {
		log.Fatalf("section=%v, iniFile=%v", section, cfg)
		return
	}
	if err := cfg.Section(section).MapTo(v); err != nil {
		panic(err)
	}
}
