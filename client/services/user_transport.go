package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func EncodeRequestFunc(ctx context.Context, req *http.Request, i interface{}) error {
	var userReq = i.(UserReq)
	req.URL.Path += fmt.Sprintf("/user/%d", userReq.Uid)
	return nil
}

func DecodeRequestFunc(ctx context.Context, resp *http.Response) (response interface{}, err error) {
	if resp.StatusCode > 400 {
		return nil, fmt.Errorf("http fail: statusCode=[%d]", resp.StatusCode)
	}
	var userResp = new(UserResp)
	if err = json.NewDecoder(resp.Body).Decode(userResp); err != nil {
		return nil, err
	}
	return userResp, nil
}
