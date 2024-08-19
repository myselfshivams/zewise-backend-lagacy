/*
Package serializers - NekoBlog backend server data serialization.
This file is for response serialization.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package serializers

type ResponseCode uint64

// BasicResponse 基本响应结构。
type BasicResponse struct {
	Code    ResponseCode `json:"code"`    // 响应码
	Message string       `json:"message"` // 响应消息
}

// DataResponse 带数据的响应结构。
type DataResponse struct {
	Code    ResponseCode `json:"code"`    // 响应码
	Message string       `json:"message"` // 响应消息
	Data    interface{}  `json:"data"`    // 响应数据
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
func NewResponse(code ResponseCode, message string, data ...interface{}) interface{} {
	// 如果没有数据，则返回基本响应。
	if len(data) == 0 {
		return BasicResponse{
			Code:    code,
			Message: message,
		}
	}

	// 如果只有一个数据，则返回带数据的响应。
	if len(data) == 1 {
		return DataResponse{
			Code:    code,
			Message: message,
			Data:    data[0],
		}
	}

	// 如果有多个数据，则返回带数据的响应。
	return DataResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
