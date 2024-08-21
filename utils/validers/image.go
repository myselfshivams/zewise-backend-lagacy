/*
Package validers - NekoBlog backend server data validation.
This file is for postimage file validation.
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

	"github.com/Kirisakiii/neko-micro-blog-backend/types"
)

// ValidImageFile 校验图片文件
//
// 参数
//   - fileHeader：文件头
//   - file：文件指针
//   - minWidth：最小宽度
//   - minHeight：最小高度
//   - maxSize：最大文件体积
//
// 返回值
//   - types.PostImageFileType：头像文件类型
//   - error：如果文件不合法，则返回相应的错误信息，否则返回nil
func ValidImageFile(fileHeader *multipart.FileHeader, file *multipart.File, minWidth, minHeight int, maxSize int64) (types.ImageFileType, error) {
	fileType := types.IMAGE_FILE_TYPE_UNKNOWN

	// 检验文件大小
	if fileHeader.Size > maxSize {
		return fileType, errors.New("image file size too large")
	}

	var (
		imgConfig image.Config
		err       error
	)

	// 校验文件类型并解码图片
	switch fileHeader.Header.Get("Content-Type") {
	case "image/jpeg":
		fileType = types.IMAGE_FILE_TYPE_JPEG
		imgConfig, err = jpeg.DecodeConfig(*file)
	case "image/png":
		fileType = types.IMAGE_FILE_TYPE_PNG
		imgConfig, err = png.DecodeConfig(*file)
	case "image/webp":
		fileType = types.IMAGE_FILE_TYPE_WEBP
		imgConfig, err = webp.DecodeConfig(*file)
	default:
		return fileType, errors.New("image file type not supported")
	}
	if err != nil {
		return fileType, err
	}

	// 校验图片尺寸
	if imgConfig.Width < minWidth || imgConfig.Height < minHeight {
		return fileType, errors.New("image size too small")
	}

	// 重置文件指针
	_, err = (*file).Seek(0, 0)
	if err != nil {
		return fileType, err
	}

	return fileType, nil
}
