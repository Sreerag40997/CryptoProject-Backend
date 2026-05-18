package rbac

import "gorm.io/gorm"

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

// Create or Delete Permission definitions
func (r *Repository) CreatePermission(p *Permission) error {
    return r.db.Create(p).Error
}

func (r *Repository) CreateRole(role *Role) error {
    return r.db.Create(role).Error
}

func (r *Repository) DeletePermission(id uint) error {
    return r.db.Delete(&Permission{}, id).Error
}

// Toggle Logic: Using GORM's Association feature
func (r *Repository) UpdateRolePermission(roleID uint, permID uint, enabled bool) error {
    var role Role
    if err := r.db.First(&role, roleID).Error; err != nil {
        return err
    }

    var perm Permission
    if err := r.db.First(&perm, permID).Error; err != nil {
        return err
    }

    if enabled {
        // Adds entry to role_permissions junction table
        return r.db.Model(&role).Association("Permissions").Append(&perm)
    } else {
        // Removes entry from role_permissions junction table
        return r.db.Model(&role).Association("Permissions").Delete(&perm)
    }
}


func (r *Repository) GetPermissionsByRole(roleName string) (*Role, error) {

	var role Role

	err := r.db.
		Preload("Permissions").
		Where("name = ?", roleName).
		First(&role).Error

	return &role, err
}

func (r *Repository) GetAllRoles() ([]Role, error) {
	var roles []Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *Repository) GetAllPermissions() ([]Permission, error) {
	var permissions []Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}