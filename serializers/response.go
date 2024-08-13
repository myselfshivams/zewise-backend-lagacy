package serializers

// BasicResponse 基本响应结构。
type BasicResponse struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

// DataResponse 带数据的响应结构。
type DataResponse struct {
	Code    uint64      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewResponse 创建一个新的响应。
//
// 参数：
//   - code：响应码
//   - message：响应消息
//   - data：响应数据
//
// 返回值：
//   - interface{}：新的响应结构体。
func NewResponse(code uint64, message string, data ...interface{}) interface{} {
	if len(data) == 0 {
		return BasicResponse{
			Code:    code,
			Message: message,
		}
	}

	if len(data) == 1 {
		return DataResponse{
			Code:    code,
			Message: message,
			Data:    data[0],
		}
	}

	return DataResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
