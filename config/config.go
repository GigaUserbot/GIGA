package config

import (
	"encoding/json"
	"os"

	"github.com/anonyindian/logger"
)

const DEBUG = false

var ValueOf = &config{}

type config struct {
	AppId         int    `json:"app_id"`
	ApiHash       string `json:"api_hash"`
	RedisUri      string `json:"redis_uri"`
	RedisCloudUrl string `json:"rediscloud_url"`
	RedisPass     string `json:"redis_pass"`
	SessionString string `json:"session_string"`
	TestServer    bool   `json:"test_mode,omitempty"`
	BotToken      string `json:"bot_token,omitempty"`
}

func Load(l *logger.Logger) {
	l = l.Create("CONFIG")
	defer l.ChangeLevel(logger.LevelMain).Println("LOADED")
	b, err := os.ReadFile("config.json")
	if err != nil {
		if err = ValueOf.setupEnvVars(); err != nil {
			l.ChangeLevel(logger.LevelError).Println(err.Error())
			os.Exit(1)
		}
		return
	}
	err = json.Unmarshal(b, ValueOf)
	if err != nil {
		l.ChangeLevel(logger.LevelError).Println("failed to load config:", err.Error())
		os.Exit(1)
	}
}
