package auth

type Repository interface {
	Create(req interface{}) error
	FindOne(model interface{}, query string, args ...any) error
}