/*
Package converters - NekoBlog backend server data converters.
This file is for avatar file converter.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package converters

import (
	"errors"
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

// ResizePostImage 调整图片大小。
//
// 参数：
//   - fileType：图片文件类型。
//   - file：图片文件
//
// 返回值：
//   - []byte：调整后的头像文件的字节切片。
//   - error：如果在调整过程中发生错误，则返回相应的错误信息，否则返回nil。
func ResizePostImage(fileType types.ImageFileType, file *multipart.File) ([]byte, error) {
	// 解码图片
	var (
		img       image.Image
		imgConfig image.Config
		err       error
	)
	switch fileType {
	case types.IMAGE_FILE_TYPE_WEBP:
		img, err = webp.Decode(*file)
		if err != nil {
			return nil, err
		}
		// 重置文件指针
		_, err = (*file).Seek(0, 0)
		if err != nil {
			return nil, err
		}
		imgConfig, err = webp.DecodeConfig(*file)

	case types.IMAGE_FILE_TYPE_JPEG:
		img, err = jpeg.Decode(*file)
		if err != nil {
			return nil, err
		}
		// 重置文件指针
		_, err = (*file).Seek(0, 0)
		if err != nil {
			return nil, err
		}
		imgConfig, err = jpeg.DecodeConfig(*file)

	case types.IMAGE_FILE_TYPE_PNG:
		img, err = png.Decode(*file)
		if err != nil {
			return nil, err
		}
		// 重置文件指针
		_, err = (*file).Seek(0, 0)
		if err != nil {
			return nil, err
		}
		imgConfig, err = png.DecodeConfig(*file)
	}
	if err != nil {
		return nil, err
	}

	// 判断图片方向并调整大小
	var resizedImg image.Image
	if imgConfig.Width > imgConfig.Height && imgConfig.Width > consts.POST_IMAGE_WIDTH_THRESHOLD {
		resizedImg = resize.Resize(
			consts.POST_IMAGE_WIDTH_THRESHOLD,
			uint(imgConfig.Height)*consts.POST_IMAGE_WIDTH_THRESHOLD/uint(imgConfig.Width),
			img,
			resize.Lanczos3,
		)
	} else if imgConfig.Width <= imgConfig.Height && imgConfig.Height > consts.POST_IMAGE_HEIGHT_THRESHOLD {
		resizedImg = resize.Resize(
			uint(imgConfig.Width)*consts.POST_IMAGE_HEIGHT_THRESHOLD/uint(imgConfig.Height),
			consts.POST_IMAGE_HEIGHT_THRESHOLD,
			img,
			resize.Lanczos3,
		)
	} else {
		resizedImg = img
	}

	if resizedImg == nil {
		return nil, errors.New("resize image failed")
	}

	// 编码图片 编码为webp存储
	imgData, err := webpEncoder.EncodeRGBA(resizedImg, consts.POST_IMAGE_QUALITY)
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
