package sql

import "github.com/anonyindian/gotgproto"

type AFK struct {
	Id     int64 `gorm:"primary_key"`
	Toggle bool
	Reason string
}

func UpdateAFK(toggle bool, reason string) {
	w := &AFK{Id: gotgproto.Self.ID}
	tx := SESSION.Begin()
	tx.FirstOrCreate(w)
	w.Toggle = toggle
	w.Reason = reason
	tx.Save(w)
	tx.Commit()
}

func GetAFK() *AFK {
	s := AFK{}
	SESSION.Where("id = ?", gotgproto.Self.GetID()).Find(&s)
	return &s
}
