package util

import jsoniter "github.com/json-iterator/go"

type CommonErr struct {
	Code  int               `json:"code"`
	Msg   string            `json:"msg"`
	Extra map[string]string `json:"extra"`
}

func (err *CommonErr) Error() string {
	var str, _ = jsoniter.MarshalToString(err)
	return str
}

func GenCommonErr(code int, msg string, extra map[string]string) error {
	return &CommonErr{
		Code:  code,
		Msg:   msg,
		Extra: extra,
	}
}
