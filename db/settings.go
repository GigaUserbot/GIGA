package db

type Settings struct {
	LogsGroup int64
	Token     string
}

type Getchatid struct {
	chatId int64
}

func UpdateSettings(logs int64, token string) {
	set("settings", &Settings{LogsGroup: logs, Token: token})
}

func Setdotsave(ChatId int64) {
	set("dotsave", &Getchatid{chatId: ChatId})
}

func Getdotsave() {
	get("dotsave", &Getchatid{})
}

func UpdateLogs(logs int64) {
	var s = &Settings{}
	get("settings", s)
	s.LogsGroup = logs
	set("settings", s)
}

func UpdateBot(token string) {
	var s = &Settings{}
	get("settings", s)
	s.Token = token
	set("settings", s)
}

func GetSettings() *Settings {
	var s = &Settings{}
	get("settings", s)
	return s
}
