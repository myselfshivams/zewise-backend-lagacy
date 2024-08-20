/*
Package type - NekoBlog backend server types.
This file is for avatar file types.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package types

// ImageFileType 图像文件类型
type ImageFileType int

// 头像文件类型
const (
	// IMAGE_FILE_TYPE_UNKNOWN 未知头像文件类型
	IMAGE_FILE_TYPE_UNKNOWN ImageFileType = -1

	// IMAGE_FILE_TYPE_WEBP WebP 格式的头像文件
	IMAGE_FILE_TYPE_WEBP ImageFileType = iota

	// IMAGE_FILE_TYPE_JPEG JPEG 格式的头像文件
	IMAGE_FILE_TYPE_JPEG

	// IMAGE_FILE_TYPE_PNG PNG 格式的头像文件
	IMAGE_FILE_TYPE_PNG
)
