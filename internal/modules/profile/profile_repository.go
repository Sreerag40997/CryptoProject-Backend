package profile

import "gorm.io/gorm"

type PgSQLRepository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &PgSQLRepository{DB: db}
}

func (r *PgSQLRepository) FindOne(model interface{}, query string, args ...any) error {
	return r.DB.Where(query, args...).First(model).Error
}

//update
func (r *PgSQLRepository) Update(model interface{}, fields map[string]interface{}, query string , args ...any) error {
	return r.DB.Model(model).Where(query, args ...).Updates(fields).Error
}

//Delete
func (r *PgSQLRepository) Delete(model interface{}, args ...any) error {
	return r.DB.Unscoped().Delete(model, args...).Error
}