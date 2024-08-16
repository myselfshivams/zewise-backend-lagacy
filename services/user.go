package services

import (
	"errors"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/converter"
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

	// 检查令牌是否达到上限
	avaliableTokens, err := service.userStore.GetUserAvaliableTokens(userAuthInfo.UserName)
	if err != nil {
		return "", err
	}
	if len(avaliableTokens) >= consts.MAX_TOKENS_PER_USER {
		// 清除最早的令牌
		err = service.userStore.BanToken(avaliableTokens[0].Token)
		if err != nil {
			return "", err
		}
	}

	// 生成令牌
	token, claims, err := generator.GenerateToken(userAuthInfo.UID, username)
	if err != nil {
		userLoginLog.Reason = "token generation error"
		inner_err := service.userStore.CreateUserLoginLog(userLoginLog)
		if inner_err != nil {
			return "", errors.Join(err, inner_err)
		}
		return "", err
	}

	// 创建可用令牌
	err = service.userStore.CreateAvaliableToken(token, claims)
	if err != nil {
		userLoginLog.Reason = "token creation error"
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

// UserUploadAvatar 用户上传头像。
//
// 参数：
//   - uid：用户ID
//   - file：头像文件
//
// 返回值：
//   - error：如果在上传过程中发生错误，则返回相应的错误信息，否则返回nil。
func (service *UserService) UserUploadAvatar(uid uint64, fileHeader *multipart.FileHeader) error {
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// 校验头像
	fileType, err := valider.ValidAvatarFile(fileHeader, &file)
	if err != nil {
		return err
	}

	// 缩放头像
	resizedAvatar, err := converter.ResizeAvatar(fileType, &file, consts.STANDERED_AVATAR_SIZE)
	if err != nil {
		return err
	}

	// 保存头像
	var sb strings.Builder
	sb.WriteString(strconv.FormatUint(uid, 10))
	sb.WriteString(".webp")
	return service.userStore.SaveAvatar(uid, sb.String(), resizedAvatar)
}

//	UserUpdatePassword 修改密码
//
//	参数：
//	- username: 用户名
//	- password: 密码
//	- newPassword: 新的密码
//
// 返回值：
//   - error：如果在上传过程中发生错误，则返回相应的错误信息，否则返回nil。
func (service *UserService) UserUpdatePassword(username string, password string, newPassword string) error {
	// 获取用户认证信息
	userAuthInfo, err := service.userStore.GetUserAuthInfo(username)
	if err != nil {
		return err
	}

	// 验证用户密码
	err = encryptor.CompareHashPassword(userAuthInfo.PasswordHash, password, userAuthInfo.Salt)
	if err != nil {
		// 密码验证失败，返回错误
		return errors.New("incorrect password")
	}

	// 取新密码哈希
	hashedNewPassword, err := encryptor.HashPassword(newPassword, userAuthInfo.Salt)
	if err != nil {
		return err
	}

	//更新密码
	err = service.userStore.UpdatePassword(userAuthInfo.UserName, hashedNewPassword)
    if err!= nil {
        return err
    }

	return nil
}


func (service *UserService) UpdateUserInfo(uid uint64, reqBody *types.UserUpdateProfile) error {
	// 构造更新Profile结构体
	updatedProfile := &models.UserInfo{}
	if reqBody.NickName != nil {
		updatedProfile.NickName = reqBody.NickName
	}
	if reqBody.Birth != nil {
		birth := time.Unix(int64(*reqBody.Birth), 0)
		updatedProfile.Birth = &birth
	}
	if reqBody.Gender != nil {
		if *reqBody.Gender != "male" && *reqBody.Gender != "female" {
			updatedProfile.Gender = reqBody.Gender
		} else {
			updatedProfile.Gender = reqBody.Gender
		}
    }

    // 在这里执行实际的数据库更新操作
    err := service.userStore.UpdateUserInfo(uid, updatedProfile)
    if err != nil {
        return err
    }

    return nil
}
