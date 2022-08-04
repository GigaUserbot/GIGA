package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/anonyindian/logger"
)

const DEBUG = false

var ValueOf = &config{}

type config struct {
	AppId         int    `json:"app_id"`
	ApiHash       string `json:"api_hash"`
	RedisUri      string `json:"redis_uri"`
	RedisPass     string `json:"redis_pass"`
	SessionString string `json:"session_string"`
	TestServer    bool   `json:"test_mode,omitempty"`
}

func Load(l *logger.Logger) {
	l = l.Create("CONFIG")
	defer l.ChangeLevel(logger.LevelMain).Println("LOADED")
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		if err := ValueOf.setupEnvVars(); err != nil {
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
