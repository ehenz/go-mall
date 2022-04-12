package forms

type SmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Type   int    `form:"type" json:"type" binding:"required,oneof=1 2"` // 1-登陆 2-注册
}
