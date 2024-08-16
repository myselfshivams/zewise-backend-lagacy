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

// CreateAvaliableToken 创建一个可用的 Token。
//
// 参数：
//   - token：Token
//   - claims：Token 的声明
//
// 返回值：
//   - error：如果在创建过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) CreateAvaliableToken(token string, claims *types.BearerTokenClaims) error {
	userAvaliableToken := &models.UserAvaliableToken{
		UID:        claims.UID,
		Username:   claims.Username,
		Token:      token,
		ExpireTime: claims.ExpiresAt.Time,
	}

	result := store.db.Create(userAvaliableToken)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// BanToken 将 Token 禁用。
//
// 参数：
//   - token：Token
//
// 返回值：
//   - error：如果在禁用过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) BanToken(token string) error {
	// 使用硬删除
	result := store.db.Where("token = ?", token).Unscoped().Delete(&models.UserAvaliableToken{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// IsTokenAvaliable 检查 Token 是否可用。
//
// 参数：
//   - token：Token
//
// 返回值：
//   - bool：如果 Token 可用，则返回 true，否则返回 false。
//   - error：如果在检查过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) IsTokenAvaliable(token string) (bool, error) {
	userAvaliableToken := new(models.UserAvaliableToken)
	result := store.db.Where("token = ?", token).First(userAvaliableToken)

	// 如果记录不存在，则返回 false
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	// 如果发生其他错误，则返回错误信息
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

// GetUserAvaliableTokens 获取用户可用的 Token。
//	
// 参数：
//   - username：用户名
//
// 返回值：
//   - []models.UserAvaliableToken：用户可用的 Token。
//   - error：如果在获取过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) GetUserAvaliableTokens(username string) ([]models.UserAvaliableToken, error) {
	tokens := make([]models.UserAvaliableToken, 0)
	// 按创建时间排序
	result := store.db.Where("username = ?", username).Order("created_at asc").Find(&tokens)
	if result.Error != nil {
		return nil, result.Error
	}
	return tokens, nil
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

// UpdatePassword 更新用户密码。
//
// 参数：
//   - username：用户名
//   - hashedNewPassword：经过哈希处理的新密码
//
// 返回值：
//   - error：如果在更新过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) UpdatePassword(username string, hashedNewPassword string) error {
	userAuthInfo := new(models.UserAuthInfo)
    result := store.db.Where("username = ?", username).First(userAuthInfo)
    if result.Error!= nil {
        return result.Error
    }

    userAuthInfo.PasswordHash = hashedNewPassword
    result = store.db.Save(userAuthInfo)
    if result.Error!= nil {
        return result.Error
    }

    return nil
}

// UpdateUserInfo 更新用户信息。
//
// 参数：
//   - uid：用户ID
//   - updatedProfile：更新后的用户信息
//
// 返回值：
//   - error：如果在更新过程中发生错误，则返回相应的错误信息，否则返回nil。
func (store *UserStore) UpdateUserInfo(uid uint64, updatedProfile *models.UserInfo) error {
	result := store.db.Model(updatedProfile).Where("id = ?", uid).Updates(updatedProfile)
    return result.Error
}