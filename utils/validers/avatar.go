/*
Package validers - NekoBlog backend server data validation.
This file is for avatar file validation.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package validers

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"

	"golang.org/x/image/webp"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
)

// ValidAvatarFile 验证头像文件
//
// 参数
//   - fileHeader：文件头
//
// 返回值
//   - types.AvatarFileType：头像文件类型
//   - error：如果文件不合法，则返回相应的错误信息，否则返回nil
func ValidAvatarFile(fileHeader *multipart.FileHeader, file *multipart.File) (types.AvatarFileType, error) {
	fileType := types.AVATAR_FILE_TYPE_UNKNOWN

	// 检验文件大小
	if fileHeader.Size > consts.MAX_AVATAR_FILE_SIZE {
		return fileType, errors.New("avatar size too large")
	}

	var (
		imgConfig image.Config
		err       error
	)

	// 校验文件类型并解码图片
	switch fileHeader.Header.Get("Content-Type") {
	case "image/jpeg":
		fileType = types.AVATAR_FILE_TYPE_JPEG
		imgConfig, err = jpeg.DecodeConfig(*file)
	case "image/png":
		fileType = types.AVATAR_FILE_TYPE_PNG
		imgConfig, err = png.DecodeConfig(*file)
	case "image/webp":
		fileType = types.AVATAR_FILE_TYPE_WEBP
		imgConfig, err = webp.DecodeConfig(*file)
	default:
		return fileType, errors.New("avatar file type not supported")
	}
	if err != nil {
		return fileType, err
	}

	// 校验图片尺寸
	if imgConfig.Width < consts.MIN_AVATAR_SIZE {
		return fileType, errors.New("avatar size too small")
	}

	// 重置文件指针
	_, err = (*file).Seek(0, 0)
	if err != nil {
		return fileType, err
	}

	return fileType, nil
}
