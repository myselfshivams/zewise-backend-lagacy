package stores

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
)

// UserStore 用户信息数据库
type UserStore struct {
	db *gorm.DB
}

// NewUserStore 返回一个新的 UserStore 实例。
//
// 返回值：
//   - *UserStore：新的 UserStore 实例。
func (factory *Factory) NewUserStore() *UserStore {
	return &UserStore{factory.db}
}

// RegisterUser 注册用户将提供的用户名、盐和哈希密码注册到数据库中。
//
// 参数：
//   - username：用户名
//   - salt：盐值
//   - hashedPassword：哈希密码
//
// 返回值：
//   - error：如果在注册过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) RegisterUser(username string, salt string, hashedPassword string) error {
	user := models.UserInfo{
		UserName: username,
		NickName: &username,
	}
	result := store.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	uid := user.ID
	userAuthInfo := models.UserAuthInfo{
		UID:          uint64(uid),
		UserName:     username,
		Salt:         salt,
		PasswordHash: hashedPassword,
	}
	result = store.db.Create(&userAuthInfo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetUserByUID 通过用户ID获取用户信息。
//
// 参数：
//   - uid：用户ID
//
// 返回值：
//   - *models.UserInfo：如果找到了相应的用户信息，则返回该用户信息，否则返回nil。
//   - error：如果在获取过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) GetUserByUID(uid uint64) (*models.UserInfo, error) {
	user := new(models.UserInfo)
	result := store.db.Where("id = ?", uid).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetUserByUsername 通过用户名获取用户信息。
//
// 参数：
//   - username：用户名
//
// 返回值：
//   - *models.UserInfo：如果找到了相应的用户信息，则返回该用户信息，否则返回nil。
//   - error：如果在获取过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) GetUserByUsername(username string) (*models.UserInfo, error) {
	user := new(models.UserInfo)
	result := store.db.Where("username = ?", username).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetUserAuthInfo 通过用户名获取用户的认证信息。
//
// 参数：
//   - username：用户名
//
// 返回值：
//   - *models.UserAuthInfo：如果找到了相应的用户认证信息，则返回该用户认证信息，否则返回nil。
//   - error：如果在获取过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) GetUserAuthInfo(username string) (*models.UserAuthInfo, error) {
	userAuthInfo := new(models.UserAuthInfo)
	result := store.db.Where("username = ?", username).First(userAuthInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return userAuthInfo, nil
}

// InsertUserLoginLog 插入用户登录日志。
//
// 参数：
//   - userLoginLogInfo：用户登录日志信息
//
// 返回值：
//   - error：如果在插入过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) CreateUserLoginLog(userLoginLogInfo *models.UserLoginLog) error {
	result := store.db.Create(userLoginLogInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// BanToken 将用户 Token 加入黑名单。
//
// 参数：
//   - token：Token
//   - claims：Token 中的声明
//
// 返回值：
//   - error：如果在插入过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) BanToken(token string, claims *types.BearerTokenClaims) error {
	userBannedToken := &models.UserBannedToken{
		UID:        claims.UID,
		Username:   claims.Username,
		Token:      token,
		ExpireTime: claims.ExpiresAt.Time,
	}

	result := store.db.Create(userBannedToken)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// IsTokenBanned 检查 Token 是否在黑名单中。
//
// 参数：
//   - token：Token
//
// 返回值：
//   - bool：如果 Token 在黑名单中，则返回true，否则返回false。
//   - error：如果在查询过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) IsTokenBanned(token string) (bool, error) {
	userBannedToken := new(models.UserBannedToken)
	result := store.db.Where("token = ?", token).First(userBannedToken)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

// SaveAvatar 保存用户头像。
//
// 参数：
//   - fileName：文件名
//   - data：文件数据
//
// 返回值：
//   - error：如果在保存过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) SaveAvatar(uid uint64, fileName string, data []byte) error {
	savePath := filepath.Join("./public/avatars", fileName)

	// 创建目标文件
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 io.Copy 将数据写入文件
	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		return err
	}

	// 获取头像文件名
	user := new(models.UserInfo)
	result := store.db.Where("id = ?", uid).First(user)
	if result.Error != nil {
		return result.Error
	}

	// 更新头像文件名
	user.Avatar = fileName
	result = store.db.Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
