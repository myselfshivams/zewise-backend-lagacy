package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/Kirisakiii/neko-micro-blog-backend/config"
	"github.com/Kirisakiii/neko-micro-blog-backend/controllers"
	"github.com/Kirisakiii/neko-micro-blog-backend/loggers"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
)

var (
	logger            *logrus.Logger
	cfg               *config.Config
	db                *gorm.DB
	controllerFactory *controllers.Factory
)

func init() {
	// logger
	logger = loggers.NewLogger()
	logger.Infoln("正在执行程序初始化...")

	var err error

	// 加载配置文件
	cfg, err = config.NewConfig()
	if err != nil {
		logger.Panicln(err.Error())
	}

	// 设置日志等级
	var (
		logLevel logrus.Level
		logMode  gormLogger.LogLevel
	)
	switch cfg.Env.Type {
	case "development":
		logLevel = logrus.DebugLevel
		logMode = gormLogger.Error
	case "production":
		logLevel = logrus.InfoLevel
		logMode = gormLogger.Silent
	default:
		logLevel = logrus.InfoLevel
		logMode = gormLogger.Silent
	}

	// 设置logrus日志等级
	logger.SetLevel(logLevel)
	logger.Debugln("日志记录等级设定为:", strings.ToUpper(logLevel.String()))

	// 连接数据库
	logger.Debugln("尝试连接至数据库...")
	db, err = gorm.Open(
		postgres.Open(fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)),
		&gorm.Config{
			Logger: gormLogger.Default.LogMode(logMode),
		},
	)
	if err != nil {
		logger.Panicln(err.Error())
	}
	logger.Debugln("数据库连接成功")

	// 迁移模型
	logger.Debugln("正在迁移数据表模型...")
	err = models.Migrate(db)
	if err != nil {
		logger.Panicln("迁移数据库模型失败：", err.Error())
	}

	// 建立控制器层工厂
	controllerFactory = controllers.NewFactory(
		services.NewFactory(
			stores.NewFactory(db),
		),
	)
}

func main() {
	app := fiber.New()
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}][${latency}][${status}][${method}] ${path}\n",
	}))
	app.Use(compress.New(compress.Config{
		Level: cfg.Compress.Level,
	}))

	// 静态资源路由
	resource := app.Group("/resources")
	resource.Static("/avatar", "./public/avatars") // 头像资源路由

	// api 路由
	api := app.Group("/api")

	// User 路由
	userController := controllerFactory.NewUserController()
	user := api.Group("/user")
	user.Get("/profile", userController.NewProfileHandler())    // 查询用户信息
	user.Post("/register", userController.NewRegisterHandler()) // 用户注册
	user.Post("/login", userController.NewLoginHandler())       // 用户登录

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)))
}
