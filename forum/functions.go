package forum

func XExists(id uint, model string) bool {
	var count int64
	var db = InitDatabase()
	switch model {
	case "question":
		db.Model(&Question{}).Where("id = ?", id).Count(&count)
		return count > 0
	case "answer":
		db.Model(&Answer{}).Where("id = ?", id).Count(&count)
		return count > 0
	default:
		return false
	}
}
