package db

func UpdateAFK(toggle bool, reason string) {
	setBool("afk", toggle)
	set("afk_reason", reason)
}

func GetAFK() bool {
	return getBool("afk")
}

func GetAFKReason() string {
	return get("afk_reason").String()
}
