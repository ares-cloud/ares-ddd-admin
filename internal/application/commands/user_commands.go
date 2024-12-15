package commands

type CreateUserCommand struct {
	Username       string  `json:"username" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Password       string  `json:"password" binding:"required"`
	Phone          string  `json:"phone"`
	Email          string  `json:"email"`
	RoleIDs        []int64 `json:"role_ids"`
	InvitationCode string  `json:"invitation_code"`
}

type UpdateUserCommand struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Email    string  `json:"email"`
	FaceURL  string  `json:"faceUrl"`
	Remark   string  `json:"remark"`
	Status   int8    `json:"status"`
	RoleIDs  []int64 `json:"roleIds"`
}

type DeleteUserCommand struct {
	ID string `json:"id"`
}

type UpdateUserStatusCommand struct {
	ID     string `json:"id"`
	Status int8   `json:"status"`
}
