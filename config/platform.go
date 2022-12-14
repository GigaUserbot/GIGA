package config

import "os"

var Platform platform

type platform int

const (
	Heroku platform = iota
	Railway
	Okteto
	Local
)

func initPlatform() {
	switch {
	case checkEnv("DYNO"):
		Platform = Heroku
	case checkEnv("RAILWAY_STATIC_URL"):
		Platform = Railway
	case checkEnv("OKTETO_TOKEN"):
		Platform = Okteto
	default:
		Platform = Local
	}
}

func checkEnv(env string) bool {
	return os.Getenv(env) != ""
}
