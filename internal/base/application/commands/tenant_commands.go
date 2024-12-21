package commands

type CreateTenantCommand struct {
	Code        string            `json:"code" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	IsDefault   int8              `json:"isDefault"`
	ExpireTime  int64             `json:"expireTime"`
	AdminUser   CreateUserCommand `json:"adminUser" binding:"required"`
}

type UpdateTenantCommand struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   int8   `json:"isDefault"`
	ExpireTime  int64  `json:"expireTime"`
}

type DeleteTenantCommand struct {
	ID string `json:"id" path:"id"`
}

type AssignTenantPermissionsCommand struct {
	TenantID      string  `json:"tenantId" binding:"required"`
	PermissionIDs []int64 `json:"permissionIds" binding:"required"`
}
