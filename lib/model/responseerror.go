package model

import "strconv"

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (re ResponseError) Error() string {
	return strconv.Itoa(re.Code) + ": " + re.Message
}
