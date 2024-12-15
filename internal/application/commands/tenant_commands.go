package commands

type CreateTenantCommand struct {
	Code        string             `json:"code" binding:"required"`
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	IsDefault   int8               `json:"is_default"`
	AdminUser   *CreateUserCommand `json:"admin_user" binding:"required"`
}

type UpdateTenantCommand struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   int8   `json:"is_default"`
}

type DeleteTenantCommand struct {
	ID string `json:"id" path:"id"`
}

type AssignTenantPermissionsCommand struct {
	TenantID      string  `json:"tenant_id" binding:"required"`
	PermissionIDs []int64 `json:"permission_ids" binding:"required"`
}
