package db

func SetLogsGroup(id int64) {
	set("logs_group", id)
}

func GetLogsGroup() (i int64) {
	i, _ = get("logs_group").Int64()
	return
}

func SetBot(token string) {
	set("bot_token", token)
}
