package db

type Afk struct {
	Toggle bool
	Reason string
}

func UpdateAFK(toggle bool, reason string) {
	set("afk", &Afk{Toggle: toggle, Reason: reason})
}

func GetAFK() *Afk {
	var a = &Afk{}
	get("afk", a)
	return a
}
