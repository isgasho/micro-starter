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
func LoadConfigWithDefault(fn func(defaultConfig *Options) *Options) {

	if fn == nil {
		logger.Warnf("use default config")
	}

	//modified config
	tmpCfg := fn(Default)

	if tmpCfg != nil {
		logger.Warnf("try to use customer config failed, use default")
	}

	//val, _ := config.Get("key.subkey3")
	// if val.String("") != "Merge" {
	// 	fmt.Println("ERROR: key.subkey3 should be 'Merge' but it is:", val.String(""))
	// }
	//  get config
	dbTypeValue, err := config.Get("DBType")
	dbType := dbTypeValue.String("")
	if err != nil && dbType != "" {
		Default.DBType = dbType
	}
	logger.Infof("DBType %+v", dbType)

	redisHostValue, err := config.Get("Redis.Host")
	redisHost := redisHostValue.String("")

	if err != nil && redisHost != "" {
		Default.Redis.Host = redisHost
	}

	pubtopicValue, err := config.Get("AsyncMessage.PubTopics")
	pubtopic := pubtopicValue.StringSlice(nil)

	if pubtopic != nil && len(pubtopic) > 0 {
		Default.Pubsub.PubTopics = pubtopic
	}

	subtopicValue, err := config.Get("AsyncMessage.SubTopics")

	subtopic := subtopicValue.StringSlice(nil)
	if subtopic != nil && len(subtopic) > 0 {
		Default.Pubsub.SubTopics = subtopic
	}

	logger.Infof("Redis Host %+v", redisHost)
}
