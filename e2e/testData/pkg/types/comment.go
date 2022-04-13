package types

type Comments struct {
	ID      int `json:"id,omitempty" position:"path"`
	UserID  int `json:"user_id,omitempty"`
	Comment int `json:"comment,omitempty"`
}
