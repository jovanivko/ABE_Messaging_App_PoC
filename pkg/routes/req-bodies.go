package routes

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Position    string `json:"position" binding:"required"`
	Department  string `json:"department" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Salary      int    `json:"salary" binding:"required"`
}

type MessageReq struct {
	From  string   `json:"from" binding:"required"`
	To    []string `json:"to" binding:"required"`
	Title []byte   `json:"title" binding:"required"`
}

type FragmentReq struct {
	MsgID      int    `json:"msg_id" binding:"required"`
	Content    []byte `json:"content" binding:"required"`
	FragmentID int    `json:"fragment_id" binding:"required"`
}
