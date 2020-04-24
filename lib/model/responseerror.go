package model

type ResponseError struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
}

func (re ResponseError) Error() string {
	return re.Message
}
