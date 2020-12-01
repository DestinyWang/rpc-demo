package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DecodeUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	// http://localhost:8000/user/101
	var vars = mux.Vars(r)
	if uidStr, ok := vars["uid"]; ok {
		var uid, _ = strconv.ParseInt(uidStr, 10, 64)
		return UserReq{
			Uid:    uid,
			Method: r.Method,
		}, nil
	}
	return nil, errors.New("param err")
}

func EncodeUserResp(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}
