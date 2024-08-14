package stores

import (
	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
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
