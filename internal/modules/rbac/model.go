package rbac

import "time"

type Role struct {
	ID          uint         `gorm:"primaryKey"`
	Name        string       `gorm:"unique;not null"` // roles here
	Permissions []Permission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time
}

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Slug string `gorm:"unique;not null"` // permision name for the code use
	Name string //  permision name for the admin view
}