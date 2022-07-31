package sql

import "github.com/anonyindian/gotgproto"

type Settings struct {
	Id        int64 `gorm:"primary_key"`
	LogsGroup int64
	// AssistantBot int64
}

func UpdateSettings(logsGroup int64) {
	w := &Settings{Id: gotgproto.Self.ID}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.LogsGroup = logsGroup
	// w.AssistantBot = assistantBot
	tx.Save(w)
	tx.Commit()
}

func GetSettings() *Settings {
	s := Settings{}
	SESSION.Where("id = ?", gotgproto.Self.GetID()).Find(&s)
	return &s
}
