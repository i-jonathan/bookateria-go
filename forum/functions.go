package forum

func QuestionExists(id uint) bool {
	var count int64
	var db = InitDatabase()
	db.Model(&Question{}).Where("id = ?", id).Count(&count)

	return count > 0
}