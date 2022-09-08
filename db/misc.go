package db

const tagLoggerKey = "tag_logger"

func TagLogger(toggle bool) {
	setBool(tagLoggerKey, toggle)
}

func GetTagLogger() bool {
	return getBool(tagLoggerKey)
}
