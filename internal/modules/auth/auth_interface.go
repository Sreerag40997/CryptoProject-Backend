package auth

type Repository interface {
	Create(req interface{}) error
	Save(model interface{}) error
	FindOne(model interface{}, query string, args ...any) error
	FindAll() ([]User, error)
	Update(model interface{}, fields map[string]interface{}, query string , args ...any) error
}