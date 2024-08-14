package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/encryptor"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/generator"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/valider"
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

// GetUserInfoByUsername 根据用户名获取用户信息。
// w
// 参数：
//   - username：用户名
//
// 返回值：
//   - *models.UserInfo：用户信息模型。
func (service *UserService) GetUserInfoByUID(uid uint64) (*models.UserInfo, error) {
	user, err := service.userStore.GetUserByUID(uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserInfoByUsername 根据用户名获取用户信息。
//
// 参数：
//   - username：用户名
//
// 返回值：
//   - *models.UserInfo：用户信息模型。
func (service *UserService) GetUserInfoByUsername(username string) (*models.UserInfo, error) {
	user, err := service.userStore.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
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

// LoginUser 用户登录。
//
// 参数：
//   - username：用户名
//   - password：密码
//   - ip：登录IP
//
// 返回值：
//   - string：Bearer Token
//   - error：如果在登录过程中发生错误，则返回相应的错误信息，否则返回nil。
func (service *UserService) LoginUser(username string, password string, ip string, app string, device string) (string, error) {
	// 获取用户认证信息
	userAuthInfo, err := service.userStore.GetUserAuthInfo(username)
	if err != nil {
		return "", err
	}

	// 构造登录日志
	userLoginLog := &models.UserLoginLog{
		UID:         userAuthInfo.UID,
		LoginTime:   time.Now(),
		LoginIP:     ip,
		Application: app,
		Device:      device,
		IsSucceed:   false,
		IfChecked:   false,
	}

	// 验证密码
	err = encryptor.CompareHashPassword(userAuthInfo.PasswordHash, password, userAuthInfo.Salt)
	if err != nil {
		userLoginLog.Reason = "password error"
		inner_err := service.userStore.CreateUserLoginLog(userLoginLog)
		if inner_err != nil {
			return "", errors.Join(err, inner_err)
		}
		return "", errors.New("password error")
	}

	// 生成令牌
	token, err := generator.GenerateToken(userAuthInfo.UID, username)
	if err != nil {
		userLoginLog.Reason = "token generation error"
		inner_err := service.userStore.CreateUserLoginLog(userLoginLog)
		if inner_err != nil {
			return "", errors.Join(err, inner_err)
		}
		return "", err
	}

	// 更新登录日志
	userLoginLog.IsSucceed = true
	userLoginLog.BearerToken = token
	err = service.userStore.CreateUserLoginLog(userLoginLog)
	if err != nil {
		return "", err
	}

	return token, nil
}
