package entity

type (
	LoginEmailInput struct {
		Email    string `form:"email" json:"email" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	DataUserInput struct {
		ID       int
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Roles    string `json:"roles"`
	}

	CreateUserInput struct {
		Name     string `json:"name" form:"name"`
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
		Username string `json:"username" form:"username"`
		Address  string `json:"address" form:"address"`
		Phone    string `json:"phone" form:"phone"`
		Avatar   string
	}
)
