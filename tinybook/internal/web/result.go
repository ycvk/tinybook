package web

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
