package model

import "JetIot/util"

type Responses_t struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func GetSuccessResponses(msg string, data ...interface{}) Responses_t {
	if len(data) == 0 {
		data = append(data, nil)
	}
	t := Responses_t{
		Msg:     msg,
		Data:    data[0],
		Code:    util.ERR_NO_ERROR,
		Success: true,
	}
	return t
}

func GetFailResponses(msg string, code int) Responses_t {
	t := Responses_t{
		Msg:     msg,
		Data:    nil,
		Code:    code,
		Success: false,
	}
	return t
}
