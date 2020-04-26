package entity

type Result struct {
	Code int `json:"code"`
	Message string `json:"msg"`
	Data interface{} `json:"data"`
}
const (
	//成功
	CODE_SUCCESS int = 1
	CODE_ERROR = -1
)
func (r *Result)SetCode(code int){
	r.Code = code
}
func (r *Result)SetMessage(msg string){
	r.Message = msg
}
func (r *Result) SetData(data map[string]interface{}){
	r.Data = data
}