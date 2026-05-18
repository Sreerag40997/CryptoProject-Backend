package rbac

//for the admin creates a new permission type
type PermissionRequest struct {
    Slug string `json:"slug" validate:"required"`
    Name string `json:"name" validate:"required"`
}

//for assign permissions for a role the toggle switches
type AssignPermissionRequest struct {
    RoleID       uint `json:"role_id" validate:"required"`
    PermissionID uint `json:"permission_id" validate:"required"`
    Enabled      bool `json:"enabled"` // true = link, false = unlink
}


type CreateRoleRequest struct {
    Name string `json:"name"`
}