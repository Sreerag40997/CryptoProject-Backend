package auth

type Repository interface {
	Create(req interface{}) error
	FindOne(model interface{}, query string, args ...any) error
	Update(model interface{}, fields map[string]interface{}, query string , args ...any) error
}