package config

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func (c *config) setupEnvVars() error {
	_ = godotenv.Load()
	val := reflect.ValueOf(c).Elem()
	notFoundArr := make([]string, 0)
	for i := 0; i < val.NumField(); i++ {
		name := toUpperSnakeCase(val.Type().Field(i).Name)
		envVal := os.Getenv(name)
		tagVal := strings.TrimSpace(val.Type().Field(i).Tag.Get("json"))
		if envVal == "" && !strings.HasSuffix(tagVal, "omitempty") {
			notFoundArr = append(notFoundArr, name)
		}
		switch val.Type().Field(i).Type.Kind() {
		case reflect.Int:
			ev, _ := strconv.ParseInt(envVal, 10, 64)
			val.Field(i).SetInt(ev)
		case reflect.String:
			val.Field(i).SetString(envVal)
		case reflect.Bool:
			ev, _ := strconv.ParseBool(envVal)
			val.Field(i).SetBool(ev)
		}
	}
	if len(notFoundArr) == 0 {
		return nil
	}
	return fmt.Errorf("failed to find following ENV vars: [%s]", strings.Join(notFoundArr, ", "))
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toUpperSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}
