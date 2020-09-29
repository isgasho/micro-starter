/*
 * @Author: Edward https://github.com/crazybber
 * @Date: 2020-09-21 17:43:58
 * @Last Modified by: Eamon
 * @Last Modified time: 2020-09-22 19:22:59
 * @Description: Configuration of current service
 */

package config

import (
	"time"

	"github.com/micro-community/auth/cache"
	"github.com/micro-community/auth/db/nosql"
	"github.com/micro-community/auth/db/sql"
	"github.com/micro-community/auth/pubsub"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service/config"
)

//Options of type
type Options struct {
	DBType  string
	Host    string
	Timeout int

	// sql db config
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration

	Redis   cache.Options
	MySQL   *sql.MySqlOptions
	SQLite  *sql.SQLiteOptions
	Mongodb *nosql.MongoOptions
	Dgraph  *nosql.DgraphOptions
	Pubsub  *pubsub.Options

	TenantKey string
}

//Service configuration and register
var (
	BASE_HERF_PATH = "./"
	Cfg            *Options //User Loaded Options,if not setted ,default value will be used.
)

//Default of config
var Default = &Options{

	DBType:          "memory",
	MaxOpenConns:    2,
	MaxIdleConns:    10,
	ConnMaxLifetime: time.Duration(time.Hour),

	Redis: cache.Options{
		MasterName:     "",
		SentinelAddrs:  nil,
		Host:           "localhost",
		Password:       "",
		DB:             0,
		MaxIdle:        1,
		MaxIdleTimeout: 1,
	},
	SQLite: &sql.SQLiteOptions{
		User:     "",
		Password: "",
		Host:     "localhost",
		DBName:   "",
		Path:     "",
	},
	Mongodb: &nosql.MongoOptions{
		User:     "",
		Password: "",
		Host:     "localhost",
		Port:     27017,
		DBName:   "auth",
	},
	Dgraph: &nosql.DgraphOptions{
		User:     "",
		Password: "",
		Host:     "localhost",
		Port:     0,
		DBName:   "",
	},
	Pubsub: &pubsub.Options{
		PubTopics: nil,
		SubTopics: nil,
	},
}

//LoadConfigWithDefault Load Options With Default
func LoadConfigWithDefault(fn func(preConfig *Options) *Options) {

	if fn == nil {
		logger.Warnf("use default config")
		Cfg = Default
	}

	//modified config
	tmpCfg := fn(Cfg)

	if tmpCfg != nil {
		logger.Warnf("try to use customer config failed, use default")
		Cfg = tmpCfg
	}

	//  get config
	dbType := config.Get("DBType").String("")
	if len(dbType) > 0 {
		Cfg.DBType = dbType
	}
	logger.Infof("DBType %+v", dbType)

	redisHost := config.Get("Redis", "Host").String("")
	if len(redisHost) > 0 {
		Cfg.Redis.Host = redisHost
	}

	pubtopic := config.Get("AsyncMessage", "PubTopics").StringSlice(nil)

	if pubtopic != nil && len(pubtopic) > 0 {
		Cfg.Pubsub.PubTopics = pubtopic
	}

	subtopic := config.Get("AsyncMessage", "SubTopics").StringSlice(nil)

	if subtopic != nil && len(subtopic) > 0 {
		Cfg.Pubsub.SubTopics = subtopic
	}

	logger.Infof("Redis Host %+v", redisHost)
}
