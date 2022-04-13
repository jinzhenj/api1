package user

type RegisterReq struct {
	ID int `json:"id,omitempty" position:"body"`
}

type LoginReq struct {
	ID int `json:"id,omitempty" position:"body"`
}

type User struct {
	ID    int     `json:"id,omitempty"`
	Roles []Roles `json:"roles,omitempty"`
}
