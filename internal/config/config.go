package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	DEVELOPER    = "developer"
	HOMOLOGATION = "homologation"
	PRODUCTION   = "production"
)

type Config struct {
	PORT          string `json:"port"`
	Mode          string `json:"mode"`
	MongoDBConfig `json:"mongo_config"`
	RedisConfig   RedisDBConfig `json:"redis_config"`
	RMQConfig     RMQConfig     `json:"rmq_config"`
}

type MongoDBConfig struct {
	MDB_URI        string `json:"mdb_uri"`
	MDB_NAME       string `json:"mdb_name"`
	MDB_COLLECTION string `json:"mdb_collection"`
}

type RedisDBConfig struct {
	RDB_HOST string `json:"rdb_host"`
	RDB_PORT string `json:"rdb_port"`
	RDB_USER string `json:"rdb_user"`
	RDB_PASS string `json:"rdb_pass"`
	RDB_DB   int64  `json:"rdb_db"`
	RDB_DSN  string `json:"-"`
}

type RMQConfig struct {
	RMQ_URI                  string `json:"rmq_uri"`
	RMQ_MAXX_RECONNECT_TIMES int    `json:"rmq_maxx_reconnect_times"`
}

func NewConfig() *Config {
	conf := defaultConf()

	SRV_PORT := os.Getenv("SRV_PORT")
	if SRV_PORT != "" {
		conf.PORT = SRV_PORT
	}

	SRV_MODE := os.Getenv("SRV_MODE")
	if SRV_MODE != "" {
		conf.Mode = SRV_MODE
	}

	SRV_RDB_HOST := os.Getenv("SRV_RDB_HOST")
	if SRV_RDB_HOST != "" {
		conf.RedisConfig.RDB_HOST = SRV_RDB_HOST
	}

	SRV_RDB_PORT := os.Getenv("SRV_RDB_PORT")
	if SRV_RDB_PORT != "" {
		conf.RedisConfig.RDB_PORT = SRV_RDB_PORT
	}

	SRV_RDB_USER := os.Getenv("SRV_RDB_USER")
	if SRV_RDB_USER != "" {
		conf.RedisConfig.RDB_USER = SRV_RDB_USER
	}

	SRV_RDB_PASS := os.Getenv("SRV_R_PASS")
	if SRV_RDB_PASS != "" {
		conf.RedisConfig.RDB_PASS = SRV_RDB_PASS
	}

	SRV_RDB_DB := os.Getenv("SRV_R_DB")
	if SRV_RDB_DB != "" {
		conf.RedisConfig.RDB_DB, _ = strconv.ParseInt(SRV_RDB_DB, 10, 64)
	}

	SRV_MDB_URI := os.Getenv("SRV_MDB_URI")
	if SRV_MDB_URI != "" {
		conf.MDB_URI = SRV_MDB_URI
	}

	SRV_MDB_NAME := os.Getenv("SRV_MDB_NAME")
	if SRV_MDB_NAME != "" {
		conf.MDB_NAME = SRV_MDB_NAME
	}

	SRV_MDB_COLLECTION := os.Getenv("SRV_MDB_COLLECTION")
	if SRV_MDB_COLLECTION != "" {
		conf.MDB_COLLECTION = SRV_MDB_COLLECTION
	}

	SRV_RDB_DSN := os.Getenv("SRV_RDB_DSN")
	if SRV_RDB_DSN != "" {
		conf.RedisConfig.RDB_DSN = SRV_MDB_COLLECTION
	}

	if len(conf.RedisConfig.RDB_HOST) > 3 {

		// "redis://<user>:<pass>@localhost:6379/<db>"
		// https://redis.uptrace.dev/guide/go-redis.html#connecting-to-redis-server

		conf.RedisConfig.RDB_DSN = fmt.Sprintf("redis://%s:%s@%s:%s/%v", conf.RedisConfig.RDB_USER, conf.RedisConfig.RDB_PASS, conf.RedisConfig.RDB_HOST, conf.RedisConfig.RDB_PORT, conf.RedisConfig.RDB_DB)
	}

	SRV_RMQ_URI := os.Getenv("SRV_RMQ_URI")
	if SRV_RMQ_URI != "" {
		conf.RMQConfig.RMQ_URI = SRV_RMQ_URI
	}

	return conf
}

func defaultConf() *Config {
	default_conf := Config{
		PORT: "8080",
		MongoDBConfig: MongoDBConfig{
			MDB_URI:        "mongodb://admin:supersenha@localhost:27017/",
			MDB_NAME:       "teste_db",
			MDB_COLLECTION: "sgStore",
		},

		Mode: DEVELOPER,

		RedisConfig: RedisDBConfig{
			RDB_HOST: "localhost",
			RDB_PORT: "6379",
			RDB_DB:   0,
			RDB_DSN:  "redis://localhost:6379/0",
		},
		RMQConfig: RMQConfig{
			RMQ_URI: "amqp://admin:supersenha@localhost:5672/",
		},
	}

	return &default_conf
}
