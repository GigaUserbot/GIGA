package db

// TODO: start database

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"time"

	"github.com/anonyindian/logger"
	"github.com/gigauserbot/giga/config"
	"github.com/go-redis/redis"
)

var client *redis.Client

func Load(l *logger.Logger) {
	l = l.Create("DATABASE")
	client = redis.NewClient(&redis.Options{
		Addr:         config.ValueOf.RedisUri,
		Password:     config.ValueOf.RedisPass,
		DB:           0,
		DialTimeout:  time.Second,
		MinIdleConns: 0,
	})
	defer l.ChangeLevel(logger.LevelMain).Println("LOADED")
}

func get(key string, T interface{}) {
	b, _ := client.Get(key).Bytes()
	gob.NewDecoder(bytes.NewBuffer(b)).Decode(T)
}

func set(key string, T interface{}) {
	setRaw(key, encode(T))
}

func getRaw(key string) *redis.StringCmd {
	return client.Get(key)
}

func setRaw(key string, v interface{}) {
	client.Set(key, v, 0)
}

func setBool(key string, value bool) {
	setRaw(key, strconv.FormatBool(value))
}

func getBool(key string) bool {
	return parseBool(getRaw(key).String())
}

func encode(v interface{}) []byte {
	buf := bytes.Buffer{}
	gob.NewEncoder(&buf).Encode(v)
	return buf.Bytes()
}

func parseBool(s string) bool {
	return s == "true"
}
