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
	return r.DB.Where(query, args...).Find(model).Error
}