package utils

type JsonData struct {
	Data interface{} `json:"data"`
}

func WrapJSON(data interface{}) JsonData {
	return JsonData{Data: data}
}
