package gormcase

// GetUsersWithPreload .
func GetUsersWithPreload() ([]UserModel, error) {
	var (
		us []UserModel
	)

	// find wanted users first
	db.Model(&UserModel{}).Find(&us)
	
	// association to another fields
	err := db.Preload("Location").Preload("Careers").Find(&us).Error
	return us, err
}
