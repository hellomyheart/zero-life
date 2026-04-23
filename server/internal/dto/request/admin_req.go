package request

type AdminListUsersReq struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Search   string `form:"search"`
}

type AdminUpdateUserReq struct {
	Nickname string `json:"nickname"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

type AdminInviteUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
}

type AdminUpdateConfigurationReq struct {
	Value string `json:"value" binding:"required"`
}

type AdminTestEmailReq struct {
	Email string `json:"email" binding:"required,email"`
}
