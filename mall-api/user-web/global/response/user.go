package response

//message UserInfoResponse {
//int32 id = 1;
//string password = 2;
//string mobile = 3;
//string nickName = 4;
//uint64 birthday = 5;
//string gender = 6;
//int32 role = 7;
//}

type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"name"`
	Mobile   string `json:"mobile"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
}
