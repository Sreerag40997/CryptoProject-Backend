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

//Save
func (r *PgSQLRepository) Save(model interface{}) error {
	return r.DB.Save(model).Error
}

func (r *PgSQLRepository) FindOne(model interface{}, query string, args ...any) error {
	return r.DB.Where(query, args...).First(model).Error
}

//find all func
func (r *PgSQLRepository) FindAll() ([]User, error) {
	var model []User
	err := r.DB.Find(&model).Error
	return model, err
}

//update
func (r *PgSQLRepository) Update(model interface{}, fields map[string]interface{}, query string , args ...any) error {
	return r.DB.Model(model).Where(query, args ...).Updates(fields).Error
}