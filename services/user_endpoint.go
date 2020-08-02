package services

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

type UserReq struct {
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	Method   string `json:"method"`
}

type UserResp struct {
	Result string `json:"result"`
}

// 定义 endPoint
func GenUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		var r = req.(UserReq)
		var result string
		switch r.Method {
		case http.MethodGet: // 根据 id 获取名称
			result = userService.GetName(r.Uid)
		case http.MethodPost: // 添加用户
			if err = userService.PutUser(r.Uid, r.Username); err != nil {
				return nil, err
			}
			result = "put succ"
		case http.MethodDelete: //删除用户
			if err = userService.DelUser(r.Uid); err != nil {
				return nil, err
			}
			result = "delete succ"
		}
		return UserResp{
			Result: result,
		}, nil
	}
}
