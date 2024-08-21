/*
Package converters - NekoBlog backend server data converters.
This file is for avatar file converter.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package converters

import (
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"

	"github.com/KononK/resize"
	webpEncoder "github.com/chai2010/webp"
	"golang.org/x/image/webp"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
)

// ResizeAvatar 调整头像大小。
//
// 参数：
//   - fileType：图片文件类型。
//	 - file：图片文件
//
// 返回值：
//   - []byte：调整后的头像文件的字节切片。
//   - error：如果在调整过程中发生错误，则返回相应的错误信息，否则返回nil。
func ResizeAvatar(fileType types.ImageFileType, file *multipart.File) ([]byte, error) {
	// 解码图片
	var (
		img image.Image
		err error
	)
	switch fileType {
	case types.IMAGE_FILE_TYPE_WEBP:
		img, err = webp.Decode(*file)

	case types.IMAGE_FILE_TYPE_JPEG:
		img, err = jpeg.Decode(*file)

	case types.IMAGE_FILE_TYPE_PNG:
		img, err = png.Decode(*file)
	}
	if err != nil {
		return nil, err
	}

	// 调整图片大小
	resizedImg := resize.Resize(
		consts.STANDERED_AVATAR_SIZE,
		consts.STANDERED_AVATAR_SIZE,
		img, resize.Lanczos3,
	)

	// 编码图片 编码为webp存储
	imgData, err := webpEncoder.EncodeRGBA(resizedImg, consts.AVATAR_QUALITY)
	if err != nil {
		return nil, err
	}

	// 重置文件指针
	_, err = (*file).Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return imgData, nil
}