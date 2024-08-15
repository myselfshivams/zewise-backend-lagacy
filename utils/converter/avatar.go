package converter

import (
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"

	"github.com/KononK/resize"
	"github.com/chai2010/webp"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
)

// ResizeAvatar 调整头像大小。
//
// 参数：
//   - fileByte：头像文件的字节切片。
//   - size：调整后的头像大小。
//
// 返回值：
//   - []byte：调整后的头像文件的字节切片。
//   - error：如果在调整过程中发生错误，则返回相应的错误信息，否则返回nil。
func ResizeAvatar(fileType types.AvatarFileType, file *multipart.File, size int) ([]byte, error) {
	// 解码图片
	var (
		img image.Image
		err error
	)
	switch fileType {
	case types.AVATAR_FILE_TYPE_WEBP:
		img, err = webp.Decode(*file)

	case types.AVATAR_FILE_TYPE_JPEG:
		img, err = jpeg.Decode(*file)

	case types.AVATAR_FILE_TYPE_PNG:
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
	imgData, err := webp.EncodeRGBA(resizedImg, consts.QUALITY)
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
