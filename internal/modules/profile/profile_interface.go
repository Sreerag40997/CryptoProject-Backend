package profile

type Repository interface {
	FindOne(model interface{}, query string, args ...any) error
	Update(model interface{}, fields map[string]interface{}, query string , args ...any) error
	Delete(model interface{}, args ...any) error
}