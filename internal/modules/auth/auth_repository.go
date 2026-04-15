package auth

import "gorm.io/gorm"

type PgSQLRepository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &PgSQLRepository{DB: db}
}

//Create User Func
func (r *PgSQLRepository) Create(req interface{}) error {
	return r.DB.Create(req).Error
}

func (r *PgSQLRepository) FindOne(model interface{}, query string, args ...any) error {
	return r.DB.Where(query, args...).First(model).Error
}

//update
func (r *PgSQLRepository) Update(model interface{}, fields map[string]interface{}, query string , args ...any) error {
	return r.DB.Model(model).Where(query, args ...).Updates(fields).Error
}