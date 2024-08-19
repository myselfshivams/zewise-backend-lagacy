/*
Package consts - NekoBlog backend server constants.
This file is for response code related constants.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package consts

import "github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"

const (
	// SUCCESS 成功
	SUCCESS serializers.ResponseCode = 0

	// SERVER_ERROR 服务器错误
	SERVER_ERROR serializers.ResponseCode = 1

	// PARAMETER_ERROR 参数错误
	PARAMETER_ERROR serializers.ResponseCode = 2

	// AUTH_ERROR 认证错误
	AUTH_ERROR serializers.ResponseCode = 3

	// NETWORK_ERROR 网络错误
	NETWORK_ERROR serializers.ResponseCode = 4
)
