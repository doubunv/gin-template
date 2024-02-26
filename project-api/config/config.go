package config

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-defaults"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"path/filepath"
	"project-api/common/tool"
	"strconv"
	"sync"
	"time"
)

var (
	Server      *YamlConfig
	MysqlClient *gorm.DB
	RedisClient *redis.Client

	serverOnce      sync.Once
	MysqlClientOnce sync.Once
	redisClientOnce sync.Once
)

type (
	YamlConfig struct {
		Name       *string        `yaml:"name"`
		Env        *string        `yaml:"env"`
		Port       *int           `yaml:"port"`
		DebugLevel yamlDebugLevel `yaml:"debugLevel"`

		Redis    yamlRedis    `yaml:"redis"`
		Mysql    yamlMysql    `yaml:"mysql"`
		RabbitMq yamlRabbitMq `yaml:"rabbitMq"`
	}

	yamlRedis struct {
		Addr          string        `yaml:"addr"`
		Password      string        `yaml:"password"`
		DB            int           `yaml:"DB" default:"10"`
		PoolSize      int           `yaml:"poolSize" default:"3"`
		DialTimeout   int           `yaml:"dialTimeout" default:"10"`
		ReadTimeout   time.Duration `yaml:"readTimeout" default:"30"`
		WriterTimeout time.Duration `yaml:"writerTimeout" default:"30"`
		PoolTimeout   time.Duration `yaml:"poolTimeout" default:"30"`
	}

	yamlMysql struct {
		Uri string `yaml:"uri"`
	}

	yamlRabbitMq struct {
		Addr string `yaml:"addr"`
	}

	yamlDebugLevel struct {
		Debug string `yaml:"debug" default:"1"`
		Info  string `yaml:"info" default:"1"`
		Error string `yaml:"error" default:"1"`
	}
)

func init() {
	logInit()
	configInit()
	mysqlInit()
	redisInit()
	zapLogInit()
}

func logInit() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func configInit() *YamlConfig {
	serverOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				panic("read config yaml file failed")
			case <-time.After(time.Nanosecond):
				env := os.Getenv("environ")
				switch env {
				case "dev":
					filePath, err := filepath.Abs("config/dev.application.yaml")
					if err != nil {
						panic(err)
					}
					viper.SetConfigFile(filePath)
					gin.SetMode(gin.DebugMode)
				default:
					filePath, err := filepath.Abs("config/application.yaml")
					if err != nil {
						panic(err)
					}
					viper.SetConfigFile(filePath)
					gin.SetMode(gin.ReleaseMode)
				}

				if err := viper.ReadInConfig(); err != nil {
					panic(fmt.Errorf("Setting.Setup.Fatal error on reading config file: %s \n", err))
				}

				if err := viper.Unmarshal(&Server); err != nil {
					panic(fmt.Errorf("Setting.Setup.Fatal error on unmarshal config: %s \n", err))
				}

				log.Printf("successfully set config, path: %s", viper.ConfigFileUsed())

				return
			}
		}
	})
	defaults.SetDefaults(&Server.DebugLevel)
	return Server
}

func redisInit() *redis.Client {
	redisClientOnce.Do(func() {
		var err error

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				panic("redis client connection failed")
			case <-time.After(time.Nanosecond):
				RedisClient = redis.NewClient(&redis.Options{
					Addr:     Server.Redis.Addr,
					Password: Server.Redis.Password, // 密码
					DB:       Server.Redis.DB,       // 默认数据库
					PoolSize: Server.Redis.PoolSize, // 连接池大小
				})
				_, err = RedisClient.Ping(ctx).Result()
				if err != nil {
					panic(err)
				}
				log.Println("successfully connect redis client")

				return
			}
		}
	})

	return RedisClient
}

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{}) {
	if _, ok := args[1].(float64); ok && args[1].(float64) > 200 {
		tool.Error(nil, "SQLLog:"+fmt.Sprintf(format, args...))
		return
	}
	tool.Info(nil, "SQLLog:"+fmt.Sprintf(format, args...))
}

func mysqlInit() *gorm.DB {
	MysqlClientOnce.Do(func() {
		var err error

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				panic("mysql client connection failed")
			case <-time.After(time.Nanosecond):
				MysqlClient, err = gorm.Open(mysql.Open(Server.Mysql.Uri), &gorm.Config{
					Logger: logger.New(
						Writer{},
						logger.Config{
							SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
							LogLevel:                  logger.Info,            // Log level
							IgnoreRecordNotFoundError: false,                  // Ignore ErrRecordNotFound error for logger
							Colorful:                  false,                  // Disable color
						},
					),
				})
				if err != nil {
					panic(err)
				}
				log.Println("successfully connect mysql client")

				return
			}
		}
	})
	return MysqlClient
}

func zapLogInit() {
	workDir, _ := os.Getwd()
	lc := &tool.LogConfig{
		DebugFileName: path.Join(workDir, "logs", "debug.log"),
		InfoFileName:  path.Join(workDir, "logs", "info.log"),
		ErrorFileName: path.Join(workDir, "logs", "error.log"),
		MaxSize:       200,
		MaxAge:        2,
		MaxBackups:    3,
	}

	debug, _ := strconv.Atoi(Server.DebugLevel.Debug)
	info, _ := strconv.Atoi(Server.DebugLevel.Info)
	errorL, _ := strconv.Atoi(Server.DebugLevel.Error)
	debugLevel := tool.DebugLevel{
		Env:   *Server.Env,
		Debug: debug,
		Info:  info,
		Error: errorL,
	}
	err := tool.InitLogger(lc, &debugLevel)
	if err != nil {
		return
	}
}
