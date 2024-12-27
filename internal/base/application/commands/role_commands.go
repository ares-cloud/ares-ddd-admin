package commands

type CreateRoleCommand struct {
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Localize    string  `json:"localize"`
	Description string  `json:"description"`
	Sequence    int     `json:"sequence"`
	Type        int8    `json:"type"`
	PermIDs     []int64 `json:"permIds"`
}

type UpdateRoleCommand struct {
	ID          int64   `json:"id" binding:"required"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Sequence    int     `json:"sequence"`
	Type        int8    `json:"type"`
	Localize    string  `json:"localize"`
	Status      *int8   `json:"status"`
	PermIDs     []int64 `json:"permIds"`
}

type DeleteRoleCommand struct {
	ID int64 `json:"id" binding:"required"`
}
