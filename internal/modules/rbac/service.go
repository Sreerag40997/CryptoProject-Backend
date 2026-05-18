package rbac

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddNewPermission(req PermissionRequest) error {
	p := &Permission{Slug: req.Slug, Name: req.Name}
	return s.repo.CreatePermission(p)
}

func (s *Service) TogglePermission(req AssignPermissionRequest) error {
	
	return s.repo.UpdateRolePermission(req.RoleID, req.PermissionID, req.Enabled)
}

func (s *Service) CreateRole(req CreateRoleRequest) error {

	role := &Role{
		Name: req.Name,
	}

	return s.repo.CreateRole(role)
}

func (s *Service) GetPermissionsByRole(role string) (*Role, error) {

	r, err := s.repo.GetPermissionsByRole(role)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Service) GetAllRoles() ([]Role, error) {
	return s.repo.GetAllRoles()
}

func (s *Service) GetAllPermissions() ([]Permission, error) {
	return s.repo.GetAllPermissions()
}