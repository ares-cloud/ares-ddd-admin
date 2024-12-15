package base_info

type Success struct {
	Code int    `json:"code" example:"200"`
	Msg  string `json:"msg" example:"err msg"`
	// 错误原因，定义为业务判定错误码
	Reason string      `json:"reason" example:"success"`
	Data   interface{} `json:"data"`
}
type Swagger400Resp struct {
	Code   int    `json:"code" example:"400"`
	Msg    string `json:"msg" example:"err msg"`
	Reason string `json:"reason" example:"err_reason"`
}
type Swagger401Resp struct {
	Code   int    `json:"code" example:"401"`
	Msg    string `json:"msg" example:"err msg"`
	Reason string `json:"reason" example:"err_reason"`
}

type Swagger500Resp struct {
	Code   int    `json:"code" example:"500"`
	Msg    string `json:"msg" example:"err msg"`
	Reason string `json:"reason" example:"err_reason"`
}
