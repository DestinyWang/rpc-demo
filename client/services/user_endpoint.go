package endpoint

type UserReq struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	Method   string `json:"method"`
}

type UserResp struct {
	Result string `json:"result"`
}
