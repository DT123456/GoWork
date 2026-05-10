package models

type Role struct {
	Model
	Name        string       `gorm:"unique" json:"name"`                  // 角色名称
	Description string       `json:"description"`                          // 角色描述
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"` // 权限关联
}

type Permission struct {
	Model
	Name        string `gorm:"unique" json:"name"`        // 权限名称，如 "user:create"
	Path        string `json:"path"`                      // 请求路径
	Method      string `json:"method"`                    // 请求方法 GET/POST/PUT/DELETE
	Description string `json:"description"`              // 权限描述
}

// UserRole 用户角色关联表
type UserRole struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}
