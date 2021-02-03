package response

import (
	"JetIot/util/errorCode"
)

type Responses_t struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	EventId int         `json:"event_id"`
}

func GetSuccessResponses(msg string, data ...interface{}) Responses_t {
	if len(data) == 0 {
		data = append(data, nil)
	}
	t := Responses_t{
		Msg:     msg,
		Data:    data[0],
		Code:    errorCode.ERR_NO_ERROR,
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

func GetMqttSuccessResponses(msg string, eventId int, data ...interface{}) Responses_t {
	if len(data) == 0 {
		data = append(data, nil)
	}
	t := Responses_t{
		Msg:     msg,
		EventId: eventId,
		Data:    data[0],
		Code:    errorCode.ERR_NO_ERROR,
		Success: true,
	}
	return t
}

func GetMqttFailResponses(msg string, eventId int, code int) Responses_t {
	t := Responses_t{
		Msg:     msg,
		Data:    nil,
		Code:    code,
		EventId: eventId,
		Success: false,
	}
	return t
}
