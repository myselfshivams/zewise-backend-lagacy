/*
Package type - NekoBlog backend server types.
This file is for avatar file types.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package types

// AvatarFileType 头像文件类型
type AvatarFileType int

// 头像文件类型
const (
	// AVATAR_FILE_TYPE_UNKNOWN 未知头像文件类型
	AVATAR_FILE_TYPE_UNKNOWN AvatarFileType = -1

	// AVATAR_FILE_TYPE_WEBP WebP 格式的头像文件
	AVATAR_FILE_TYPE_WEBP AvatarFileType = iota

	// AVATAR_FILE_TYPE_JPEG JPEG 格式的头像文件
	AVATAR_FILE_TYPE_JPEG
	
	// AVATAR_FILE_TYPE_PNG PNG 格式的头像文件
	AVATAR_FILE_TYPE_PNG
)
