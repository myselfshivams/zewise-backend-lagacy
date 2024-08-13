package services

import (
	"errors"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/encryptor"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/generator"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/valider"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	userStore *stores.UserStore
}

// NewUserService 返回一个新的 UserService 实例。
//
// 返回值：
//   - *UserService：新的 UserService 实例。
func (factory *Factory) NewUserService() *UserService {
	return &UserService{
		userStore: factory.storeFactory.NewUserStore(),
	}
}

// RegisterUser 注册用户。
//
// 参数：
//   - username：用户名
//   - password：密码
//
// 返回值：
//   - error：如果在注册过程中发生错误，则返回相应的错误信息，否则返回nil。
func (service *UserService) RegisterUser(username string, password string) error {
	// 验证用户名和密码是否合法
	if !valider.IsValidUsername(username) {
		return errors.New("invalid username")
	}
	if !valider.IsValidPassword(password) {
		return errors.New("invalid password")
	}

	// 检验用户名是否重复
	_, err := service.userStore.GetUserByUsername(username)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("username already exists")
	}

	// 生成盐和哈希密码
	salt, err := generator.GenerateSalt(consts.SALT_LENGTH)
	if err != nil {
		return err
	}
	hashedPassword, err := encryptor.HashPassword(password, salt)
	if err != nil {
		return err
	}

	// 注册用户
	err = service.userStore.RegisterUser(username, salt, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
