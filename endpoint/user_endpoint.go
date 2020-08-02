package endpoint

import (
	"context"
	"github.com/DestinyWang/gokit-test/services"
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
func GenUserEndpoint(userService services.IUserService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		var r = req.(UserReq)
		var result string
		switch r.Method {
		case http.MethodGet: // 根据 id 获取名称
			result = userService.GetName(r.Uid)
		case http.MethodPost: // 添加用户
			result = userService.PutUser(r.Uid, r.Username).Error()
		case http.MethodDelete: //删除用户
			result = userService.DelUser(r.Uid).Error()
		}
		return UserResp{
			Result: result,
		}, nil
	}
}
