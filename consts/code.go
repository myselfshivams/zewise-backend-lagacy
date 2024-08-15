package consts

import "github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"

const (
	// Success 成功
	SUCCESS serializers.ResponseCode = 0

	// SERVER_ERROR 服务器错误
	SERVER_ERROR serializers.ResponseCode = 1

	// PARAMETER_ERROR 参数错误
	PARAMETER_ERROR serializers.ResponseCode = 2

	// AUTH_ERROR 认证错误
	AUTH_ERROR serializers.ResponseCode = 3
)
