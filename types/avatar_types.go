package types

type AvatarFileType int

const (
	AVATAR_FILE_TYPE_UNKNOWN AvatarFileType = -1
	AVATAR_FILE_TYPE_WEBP    AvatarFileType = iota
	AVATAR_FILE_TYPE_JPEG
	AVATAR_FILE_TYPE_PNG
)
